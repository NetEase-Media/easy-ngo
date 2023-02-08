// Copyright 2022 NetEase Media Technology（Beijing）Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xredis

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"sync"
	"time"
	_ "unsafe"

	"github.com/hashicorp/go-multierror"

	"github.com/go-redis/redis/v8"
)

const (
	shardedFactor = 160
)

var keyTagPattern = regexp.MustCompile(`{(.+?)}`)

var pipelineCmds = [...]string{
	// key
	"del", "exists", "expire", "expireat", "persist", "pexpire", "pexpireat", "pttl", "randomkey", "rename", "renamenx", "restore", "sort", "ttl", "type",
	// string
	"append", "bitcount", "bitop", "decr", "decrby", "get", "getbit", "getrange", "getset", "incr", "incrby", "incrbyfloat", "mget", "mset", "msetnx", "psetex", "set", "setbit", "setex", "setnx", "setrange", "strlen",
	// hash
	"hdel", "hexists", "hget", "hgetall", "hincrby", "hincrbyfloat", "hkeys", "hlen", "hmget", "hmset", "hset", "hsetnx", "hvals", "hscan",
	// list
	"blpop", "brpop", "brpoplpush", "lindex", "linsert", "llen", "lpop", "lpush", "lpushx", "lrange", "lrem", "lset", "ltrim", "rpop", "rpoplpush", "rpush", "rpushx",
	// set
	"sadd", "scard", "sdiff", "sdiffstore", "sinter", "sinterstore", "sismember", "smembers", "smove", "spop", "srandmember", "srem", "sunion", "sunionstore", "sscan",
	// sort set
	"zadd", "zcard", "zcount", "zincrby", "zrange", "zrangebyscore", "zrank", "zrem", "zremrangebyrank", "zremrangebyscore", "zrevrange", "zrevrangebyscore", "zrevrank", "zscore", "zunionstore", "zinterstore", "zscan",
}

func NewShardedClient(sis []*ShardInfo) Redis {
	nodes := make(map[int64]*ShardInfo, len(sis)*shardedFactor)
	sortedHashes := make([]int64, 0, len(sis)*shardedFactor)
	resources := make(map[string]*ShardInfo, len(sis))
	algo := &MurmurHash{}

	for i, si := range sis {
		if si.name == "" {
			for n := 0; n < shardedFactor*si.weight; n++ {
				hash := algo.hash(fmt.Sprintf("SHARD-%d-NODE-%d", i, n))
				nodes[hash] = si
				sortedHashes = append(sortedHashes, hash)
			}
		} else {
			for n := 0; n < shardedFactor*si.weight; n++ {
				hash := algo.hash(fmt.Sprintf("%s*%d%d", si.name, si.weight, n))
				nodes[hash] = si
				sortedHashes = append(sortedHashes, hash)
			}
		}
		resources[si.id] = si
	}
	sort.Slice(sortedHashes, func(i int, j int) bool {
		return sortedHashes[i] < sortedHashes[j]
	})
	return &ShardedClient{
		nodes:        nodes,
		sortedHashes: sortedHashes,
		algo:         algo,
		resources:    resources,
		tagPattern:   keyTagPattern,
		ctx:          context.Background(),
	}
}

type ShardedClient struct {
	Redis
	nodes        map[int64]*ShardInfo
	sortedHashes []int64
	algo         Hashing
	resources    map[string]*ShardInfo
	tagPattern   *regexp.Regexp

	sync.RWMutex // 目前只有读 没有写

	ctx context.Context
}

// --- commands --------------------------------------
func (c *ShardedClient) Close() error {
	var mulerr error
	for _, r := range c.getAllShards() {
		err := r.client.Close()
		if err != nil {
			mulerr = multierror.Append(mulerr, err)
		}
	}
	return mulerr
}

func (c *ShardedClient) Pipeline() redis.Pipeliner {
	return NewShardedPipeline(c.ctx, c.processPipeline)
}
func (c *ShardedClient) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return c.Pipeline().Pipelined(ctx, fn)
}

func (c *ShardedClient) processPipeline(ctx context.Context, cmds []redis.Cmder) error {
	err := checkCmds(cmds)
	if err != nil {
		return err
	}
	cmdsMap := newCmdsMap()
	for i := range cmds {
		key, ok := cmds[i].Args()[1].(string)
		if !ok {
			return fmt.Errorf("exists unsupport command: [%s]", cmds[i])
		}
		client := c.getShard(key)
		cmdsMap.Add(client, cmds[i])
	}
	var wg sync.WaitGroup
	for client, cmds := range cmdsMap.m {
		wg.Add(1)
		go func(client Redis, cmds []redis.Cmder) {
			defer wg.Done()
			processPipeline(ctx, client, cmds)
		}(client, cmds)
	}
	wg.Wait()
	return cmdsFirstErr(cmds)
}

func checkCmds(cmds []redis.Cmder) error {
	for i := range cmds {
		if cmds[i] == nil || len(cmds[i].Args()) < 2 {
			return fmt.Errorf("unsupport command: [%s]", cmds[i])
		}
		cmd := cmds[i].Args()[0]
		exit := false
		for j := range pipelineCmds {
			if cmd == pipelineCmds[j] {
				exit = true
				break
			}
		}
		if !exit {
			return fmt.Errorf("unsupport command: [%s]", cmd)
		}
	}
	return nil
}

func (c *ShardedClient) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	panic("unsupport method..")
}
func (c *ShardedClient) TxPipeline() redis.Pipeliner {
	panic("unsupport method..")
}

func (c *ShardedClient) processTxPipeline(context.Context, []redis.Cmder) error {
	return nil
}

func (c *ShardedClient) Command(ctx context.Context) *redis.CommandsInfoCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClientGetName(ctx context.Context) *redis.StringCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Echo(ctx context.Context, message interface{}) *redis.StringCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Ping(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Quit(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.Del(ctx, keys[0])
}
func (c *ShardedClient) Unlink(ctx context.Context, keys ...string) *redis.IntCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.Unlink(ctx, keys[0])
}
func (c *ShardedClient) Dump(ctx context.Context, key string) *redis.StringCmd {
	client := c.getShard(key)
	return client.Dump(ctx, key)
}
func (c *ShardedClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.Exists(ctx, keys[0])
}
func (c *ShardedClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	client := c.getShard(key)
	return client.Expire(ctx, key, expiration)
}
func (c *ShardedClient) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	client := c.getShard(key)
	return client.ExpireAt(ctx, key, tm)
}
func (c *ShardedClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	client := c.getShard(key)
	return client.Migrate(ctx, host, port, key, db, timeout)
}
func (c *ShardedClient) Move(ctx context.Context, key string, db int) *redis.BoolCmd {
	client := c.getShard(key)
	return client.Move(ctx, key, db)
}
func (c *ShardedClient) ObjectRefCount(ctx context.Context, key string) *redis.IntCmd {
	client := c.getShard(key)
	return client.ObjectRefCount(ctx, key)
}
func (c *ShardedClient) ObjectEncoding(ctx context.Context, key string) *redis.StringCmd {
	client := c.getShard(key)
	return client.ObjectEncoding(ctx, key)
}
func (c *ShardedClient) ObjectIdleTime(ctx context.Context, key string) *redis.DurationCmd {
	client := c.getShard(key)
	return client.ObjectIdleTime(ctx, key)
}
func (c *ShardedClient) Persist(ctx context.Context, key string) *redis.BoolCmd {
	client := c.getShard(key)
	return client.Persist(ctx, key)
}
func (c *ShardedClient) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	client := c.getShard(key)
	return client.PExpire(ctx, key, expiration)
}
func (c *ShardedClient) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	client := c.getShard(key)
	return client.PExpireAt(ctx, key, tm)
}
func (c *ShardedClient) PTTL(ctx context.Context, key string) *redis.DurationCmd {
	client := c.getShard(key)
	return client.PTTL(ctx, key)
}
func (c *ShardedClient) RandomKey(ctx context.Context) *redis.StringCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Rename(ctx context.Context, key, newkey string) *redis.StatusCmd {
	client := c.getShard(key)
	return client.Rename(ctx, key, newkey)
}
func (c *ShardedClient) RenameNX(ctx context.Context, key, newkey string) *redis.BoolCmd {
	client := c.getShard(key)
	return client.RenameNX(ctx, key, newkey)
}
func (c *ShardedClient) Restore(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	client := c.getShard(key)
	return client.Restore(ctx, key, ttl, value)
}
func (c *ShardedClient) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	client := c.getShard(key)
	return client.RestoreReplace(ctx, key, ttl, value)
}
func (c *ShardedClient) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.Sort(ctx, key, sort)
}
func (c *ShardedClient) SortStore(ctx context.Context, key, store string, sort *redis.Sort) *redis.IntCmd {
	client := c.getShard(key)
	return client.SortStore(ctx, key, store, sort)
}
func (c *ShardedClient) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) *redis.SliceCmd {
	client := c.getShard(key)
	return client.SortInterfaces(ctx, key, sort)
}
func (c *ShardedClient) Touch(ctx context.Context, keys ...string) *redis.IntCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.Touch(ctx, keys[0])
}
func (c *ShardedClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	client := c.getShard(key)
	return client.TTL(ctx, key)
}
func (c *ShardedClient) Type(ctx context.Context, key string) *redis.StatusCmd {
	client := c.getShard(key)
	return client.Type(ctx, key)
}
func (c *ShardedClient) Append(ctx context.Context, key, value string) *redis.IntCmd {
	client := c.getShard(key)
	return client.Append(ctx, key, value)
}
func (c *ShardedClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	client := c.getShard(key)
	return client.Decr(ctx, key)
}
func (c *ShardedClient) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	client := c.getShard(key)
	return client.DecrBy(ctx, key, decrement)
}
func (c *ShardedClient) Get(ctx context.Context, key string) *redis.StringCmd {
	client := c.getShard(key)
	return client.Get(ctx, key)
}
func (c *ShardedClient) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {
	client := c.getShard(key)
	return client.GetRange(ctx, key, start, end)
}
func (c *ShardedClient) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	client := c.getShard(key)
	return client.GetSet(ctx, key, value)
}
func (c *ShardedClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	client := c.getShard(key)
	return client.Incr(ctx, key)
}
func (c *ShardedClient) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	client := c.getShard(key)
	return client.IncrBy(ctx, key, value)
}
func (c *ShardedClient) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	client := c.getShard(key)
	return client.IncrByFloat(ctx, key, value)
}
func (c *ShardedClient) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.MGet(ctx, keys[0])
}
func (c *ShardedClient) MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	client := c.getShard(key)
	return client.Set(ctx, key, value, expiration)
}
func (c *ShardedClient) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	client := c.getShard(key)
	return client.SetEX(ctx, key, value, expiration)
}
func (c *ShardedClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	client := c.getShard(key)
	return client.SetNX(ctx, key, value, expiration)
}
func (c *ShardedClient) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	client := c.getShard(key)
	return client.SetXX(ctx, key, value, expiration)
}
func (c *ShardedClient) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	client := c.getShard(key)
	return client.SetRange(ctx, key, offset, value)
}
func (c *ShardedClient) StrLen(ctx context.Context, key string) *redis.IntCmd {
	client := c.getShard(key)
	return client.StrLen(ctx, key)
}
func (c *ShardedClient) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	client := c.getShard(key)
	return client.GetBit(ctx, key, offset)
}
func (c *ShardedClient) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	client := c.getShard(key)
	return client.SetBit(ctx, key, offset, value)
}
func (c *ShardedClient) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	client := c.getShard(key)
	return client.BitCount(ctx, key, bitCount)
}
func (c *ShardedClient) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.BitOpAnd(ctx, destKey, keys[0])
}
func (c *ShardedClient) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.BitOpOr(ctx, destKey, keys[0])
}
func (c *ShardedClient) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.BitOpXor(ctx, destKey, keys[0])
}
func (c *ShardedClient) BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd {
	client := c.getShard(key)
	return client.BitOpNot(ctx, destKey, key)
}
func (c *ShardedClient) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd {
	client := c.getShard(key)
	return client.BitPos(ctx, key, bit, pos...)
}
func (c *ShardedClient) BitField(ctx context.Context, key string, args ...interface{}) *redis.IntSliceCmd {
	client := c.getShard(key)
	return client.BitField(ctx, key, args...)
}
func (c *ShardedClient) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	client := c.getShard(key)
	return client.SScan(ctx, key, cursor, match, count)
}
func (c *ShardedClient) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	client := c.getShard(key)
	return client.HScan(ctx, key, cursor, match, count)
}
func (c *ShardedClient) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	client := c.getShard(key)
	return client.ZScan(ctx, key, cursor, match, count)
}
func (c *ShardedClient) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	client := c.getShard(key)
	return client.HDel(ctx, key, fields...)
}
func (c *ShardedClient) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	client := c.getShard(key)
	return client.HExists(ctx, key, field)
}
func (c *ShardedClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	client := c.getShard(key)
	return client.HGet(ctx, key, field)
}
func (c *ShardedClient) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	client := c.getShard(key)
	return client.HGetAll(ctx, key)
}
func (c *ShardedClient) HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	client := c.getShard(key)
	return client.HIncrBy(ctx, key, field, incr)
}
func (c *ShardedClient) HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {
	client := c.getShard(key)
	return client.HIncrByFloat(ctx, key, field, incr)
}
func (c *ShardedClient) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.HKeys(ctx, key)
}
func (c *ShardedClient) HLen(ctx context.Context, key string) *redis.IntCmd {
	client := c.getShard(key)
	return client.HLen(ctx, key)
}
func (c *ShardedClient) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	client := c.getShard(key)
	return client.HMGet(ctx, key, fields...)
}
func (c *ShardedClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.HSet(ctx, key, values...)
}
func (c *ShardedClient) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	client := c.getShard(key)
	return client.HMSet(ctx, key, values...)
}
func (c *ShardedClient) HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	client := c.getShard(key)
	return client.HSetNX(ctx, key, field, value)
}
func (c *ShardedClient) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.HVals(ctx, key)
}
func (c *ShardedClient) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.BLPop(ctx, timeout, keys[0])
}
func (c *ShardedClient) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.BRPop(ctx, timeout, keys[0])
}
func (c *ShardedClient) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	client := c.getShard(key)
	return client.LIndex(ctx, key, index)
}
func (c *ShardedClient) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.LInsert(ctx, key, op, pivot, value)
}
func (c *ShardedClient) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.LInsertBefore(ctx, key, pivot, value)
}
func (c *ShardedClient) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.LInsertAfter(ctx, key, pivot, value)
}
func (c *ShardedClient) LLen(ctx context.Context, key string) *redis.IntCmd {
	client := c.getShard(key)
	return client.LLen(ctx, key)
}
func (c *ShardedClient) LPop(ctx context.Context, key string) *redis.StringCmd {
	client := c.getShard(key)
	return client.LPop(ctx, key)
}
func (c *ShardedClient) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.LPush(ctx, key, values...)
}
func (c *ShardedClient) LPushX(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.LPushX(ctx, key, values...)
}
func (c *ShardedClient) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.LRange(ctx, key, start, stop)
}
func (c *ShardedClient) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.LRem(ctx, key, count, value)
}
func (c *ShardedClient) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	client := c.getShard(key)
	return client.LSet(ctx, key, index, value)
}
func (c *ShardedClient) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	client := c.getShard(key)
	return client.LTrim(ctx, key, start, stop)
}
func (c *ShardedClient) RPop(ctx context.Context, key string) *redis.StringCmd {
	client := c.getShard(key)
	return client.RPop(ctx, key)
}
func (c *ShardedClient) RPopLPush(ctx context.Context, source, destination string) *redis.StringCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.RPush(ctx, key, values...)
}
func (c *ShardedClient) RPushX(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.RPushX(ctx, key, values...)
}
func (c *ShardedClient) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.SAdd(ctx, key, members...)
}
func (c *ShardedClient) SCard(ctx context.Context, key string) *redis.IntCmd {
	client := c.getShard(key)
	return client.SCard(ctx, key)
}
func (c *ShardedClient) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.SDiff(ctx, keys[0])
}
func (c *ShardedClient) SDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.SDiffStore(ctx, destination, keys[0])
}
func (c *ShardedClient) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.SInter(ctx, keys[0])
}
func (c *ShardedClient) SInterStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.SInterStore(ctx, destination, keys[0])
}
func (c *ShardedClient) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	client := c.getShard(key)
	return client.SIsMember(ctx, key, member)
}
func (c *ShardedClient) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.SMembers(ctx, key)
}
func (c *ShardedClient) SMembersMap(ctx context.Context, key string) *redis.StringStructMapCmd {
	client := c.getShard(key)
	return client.SMembersMap(ctx, key)
}
func (c *ShardedClient) SMove(ctx context.Context, source, destination string, member interface{}) *redis.BoolCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) SPop(ctx context.Context, key string) *redis.StringCmd {
	client := c.getShard(key)
	return client.SPop(ctx, key)
}
func (c *ShardedClient) SPopN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.SPopN(ctx, key, count)
}
func (c *ShardedClient) SRandMember(ctx context.Context, key string) *redis.StringCmd {
	client := c.getShard(key)
	return client.SRandMember(ctx, key)
}
func (c *ShardedClient) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.SRandMemberN(ctx, key, count)
}
func (c *ShardedClient) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.SRem(ctx, key, members...)
}
func (c *ShardedClient) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.SUnion(ctx, keys[0])
}
func (c *ShardedClient) SUnionStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.SUnionStore(ctx, destination, keys[0])
}
func (c *ShardedClient) XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XDel(ctx context.Context, stream string, ids ...string) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XLen(ctx context.Context, stream string) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XRange(ctx context.Context, stream, start, stop string) *redis.XMessageSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XRangeN(ctx context.Context, stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XRevRange(ctx context.Context, stream string, start, stop string) *redis.XMessageSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) *redis.XMessageSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XReadStreams(ctx context.Context, streams ...string) *redis.XStreamSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XGroupCreate(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XGroupSetID(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XGroupDestroy(ctx context.Context, stream, group string) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XAck(ctx context.Context, stream, group string, ids ...string) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XPending(ctx context.Context, stream, group string) *redis.XPendingCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XClaim(ctx context.Context, a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XClaimJustID(ctx context.Context, a *redis.XClaimArgs) *redis.StringSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) XTrim(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	client := c.getShard(key)
	return client.XTrim(ctx, key, maxLen)
}
func (c *ShardedClient) XTrimApprox(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	client := c.getShard(key)
	return client.XTrimApprox(ctx, key, maxLen)
}
func (c *ShardedClient) XInfoGroups(ctx context.Context, key string) *redis.XInfoGroupsCmd {
	client := c.getShard(key)
	return client.XInfoGroups(ctx, key)
}
func (c *ShardedClient) XInfoStream(ctx context.Context, key string) *redis.XInfoStreamCmd {
	client := c.getShard(key)
	return client.XInfoStream(ctx, key)
}
func (c *ShardedClient) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.BZPopMax(ctx, timeout, keys[0])
}
func (c *ShardedClient) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.BZPopMin(ctx, timeout, keys[0])
}
func (c *ShardedClient) ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZAdd(ctx, key, members...)
}
func (c *ShardedClient) ZAddNX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZAddNX(ctx, key, members...)
}
func (c *ShardedClient) ZAddXX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZAddXX(ctx, key, members...)
}
func (c *ShardedClient) ZAddCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZAddCh(ctx, key, members...)
}
func (c *ShardedClient) ZAddNXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZAddNXCh(ctx, key, members...)
}
func (c *ShardedClient) ZAddXXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZAddXXCh(ctx, key, members...)
}
func (c *ShardedClient) ZIncr(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	client := c.getShard(key)
	return client.ZIncr(ctx, key, member)
}
func (c *ShardedClient) ZIncrNX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	client := c.getShard(key)
	return client.ZIncrNX(ctx, key, member)
}
func (c *ShardedClient) ZIncrXX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	client := c.getShard(key)
	return client.ZIncrXX(ctx, key, member)
}
func (c *ShardedClient) ZCard(ctx context.Context, key string) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZCard(ctx, key)
}
func (c *ShardedClient) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZCount(ctx, key, min, max)
}
func (c *ShardedClient) ZLexCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZLexCount(ctx, key, min, max)
}
func (c *ShardedClient) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	client := c.getShard(key)
	return client.ZIncrBy(ctx, key, increment, member)
}
func (c *ShardedClient) ZInterStore(ctx context.Context, destination string, store *redis.ZStore) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	client := c.getShard(key)
	return client.ZPopMax(ctx, key, count...)
}
func (c *ShardedClient) ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	client := c.getShard(key)
	return client.ZPopMin(ctx, key, count...)
}
func (c *ShardedClient) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.ZRange(ctx, key, start, stop)
}
func (c *ShardedClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	client := c.getShard(key)
	return client.ZRangeWithScores(ctx, key, start, stop)
}
func (c *ShardedClient) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.ZRangeByScore(ctx, key, opt)
}
func (c *ShardedClient) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.ZRangeByLex(ctx, key, opt)
}
func (c *ShardedClient) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	client := c.getShard(key)
	return client.ZRangeByScoreWithScores(ctx, key, opt)
}
func (c *ShardedClient) ZRank(ctx context.Context, key, member string) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZRank(ctx, key, member)
}
func (c *ShardedClient) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZRem(ctx, key, members...)
}
func (c *ShardedClient) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZRemRangeByRank(ctx, key, start, stop)
}
func (c *ShardedClient) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZRemRangeByScore(ctx, key, min, max)
}
func (c *ShardedClient) ZRemRangeByLex(ctx context.Context, key, min, max string) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZRemRangeByLex(ctx, key, min, max)
}
func (c *ShardedClient) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.ZRevRange(ctx, key, start, stop)
}
func (c *ShardedClient) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	client := c.getShard(key)
	return client.ZRevRangeWithScores(ctx, key, start, stop)
}
func (c *ShardedClient) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.ZRevRangeByScore(ctx, key, opt)
}
func (c *ShardedClient) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.ZRevRangeByLex(ctx, key, opt)
}
func (c *ShardedClient) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	client := c.getShard(key)
	return client.ZRevRangeByScoreWithScores(ctx, key, opt)
}
func (c *ShardedClient) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {
	client := c.getShard(key)
	return client.ZRevRank(ctx, key, member)
}
func (c *ShardedClient) ZScore(ctx context.Context, key, member string) *redis.FloatCmd {
	client := c.getShard(key)
	return client.ZScore(ctx, key, member)
}
func (c *ShardedClient) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) PFAdd(ctx context.Context, key string, els ...interface{}) *redis.IntCmd {
	client := c.getShard(key)
	return client.PFAdd(ctx, key, els...)
}
func (c *ShardedClient) PFCount(ctx context.Context, keys ...string) *redis.IntCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.PFCount(ctx, keys[0])
}
func (c *ShardedClient) PFMerge(ctx context.Context, dest string, keys ...string) *redis.StatusCmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.PFMerge(ctx, dest, keys[0])
}
func (c *ShardedClient) BgRewriteAOF(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) BgSave(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClientKill(ctx context.Context, ipPort string) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClientKillByFilter(ctx context.Context, keys ...string) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClientList(ctx context.Context) *redis.StringCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClientPause(ctx context.Context, dur time.Duration) *redis.BoolCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClientID(ctx context.Context) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ConfigGet(ctx context.Context, parameter string) *redis.SliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ConfigResetStat(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ConfigSet(ctx context.Context, parameter, value string) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ConfigRewrite(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) DBSize(ctx context.Context) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) FlushAll(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) FlushAllAsync(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) FlushDB(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) FlushDBAsync(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Info(ctx context.Context, section ...string) *redis.StringCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) LastSave(ctx context.Context) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Save(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Shutdown(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ShutdownSave(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ShutdownNoSave(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) SlaveOf(ctx context.Context, host, port string) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Time(ctx context.Context) *redis.TimeCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) DebugObject(ctx context.Context, key string) *redis.StringCmd {
	client := c.getShard(key)
	return client.DebugObject(ctx, key)
}
func (c *ShardedClient) ReadOnly(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ReadWrite(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) MemoryUsage(ctx context.Context, key string, samples ...int) *redis.IntCmd {
	client := c.getShard(key)
	return client.MemoryUsage(ctx, key, samples...)
}
func (c *ShardedClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.Eval(ctx, script, keys[:1], args)
}
func (c *ShardedClient) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	if len(keys) != 1 {
		panic("the length of keys must be equal to 1")
	}
	client := c.getShard(keys[0])
	return client.EvalSha(ctx, sha1, keys[:1], args)
}
func (c *ShardedClient) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ScriptFlush(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ScriptKill(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) PubSubChannels(ctx context.Context, pattern string) *redis.StringSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) PubSubNumSub(ctx context.Context, channels ...string) *redis.StringIntMapCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) PubSubNumPat(ctx context.Context) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterSlots(ctx context.Context) *redis.ClusterSlotsCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterNodes(ctx context.Context) *redis.StringCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterMeet(ctx context.Context, host, port string) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterForget(ctx context.Context, nodeID string) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterReplicate(ctx context.Context, nodeID string) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterResetSoft(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterResetHard(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterInfo(ctx context.Context) *redis.StringCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterKeySlot(ctx context.Context, key string) *redis.IntCmd {
	client := c.getShard(key)
	return client.ClusterKeySlot(ctx, key)
}
func (c *ShardedClient) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *redis.StringSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterCountFailureReports(ctx context.Context, nodeID string) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterCountKeysInSlot(ctx context.Context, slot int) *redis.IntCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterDelSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterDelSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterSaveConfig(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterSlaves(ctx context.Context, nodeID string) *redis.StringSliceCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterFailover(ctx context.Context) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterAddSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) ClusterAddSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {
	panic("unsupport method..")
}
func (c *ShardedClient) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	client := c.getShard(key)
	return client.GeoAdd(ctx, key, geoLocation...)
}
func (c *ShardedClient) GeoPos(ctx context.Context, key string, members ...string) *redis.GeoPosCmd {
	client := c.getShard(key)
	return client.GeoPos(ctx, key, members...)
}
func (c *ShardedClient) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	client := c.getShard(key)
	return client.GeoRadius(ctx, key, longitude, latitude, query)
}
func (c *ShardedClient) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.IntCmd {
	client := c.getShard(key)
	return client.GeoRadiusStore(ctx, key, longitude, latitude, query)
}
func (c *ShardedClient) GeoRadiusByMember(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	client := c.getShard(key)
	return client.GeoRadiusByMember(ctx, key, member, query)
}
func (c *ShardedClient) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.IntCmd {
	client := c.getShard(key)
	return client.GeoRadiusByMemberStore(ctx, key, member, query)
}
func (c *ShardedClient) GeoDist(ctx context.Context, key string, member1, member2, unit string) *redis.FloatCmd {
	client := c.getShard(key)
	return client.GeoDist(ctx, key, member1, member2, unit)
}
func (c *ShardedClient) GeoHash(ctx context.Context, key string, members ...string) *redis.StringSliceCmd {
	client := c.getShard(key)
	return client.GeoHash(ctx, key, members...)
}

// ---------------------------------------------------
func (c *ShardedClient) getAllShards() []*ShardInfo {
	sis := make([]*ShardInfo, 0, len(c.resources))
	for _, v := range c.resources {
		sis = append(sis, v)
	}
	return sis
}

func (c *ShardedClient) getShard(key string) Redis {
	info := c.getShardInfo(key)
	return info.client
}

func (c *ShardedClient) getShardInfo(key string) *ShardInfo {
	c.RLock()
	defer c.RUnlock()

	k := c.algo.hash(c.getKeyTag(key))
	idx := sort.Search(len(c.sortedHashes), func(i int) bool {
		return c.sortedHashes[i] >= k
	})

	if idx >= len(c.sortedHashes) {
		idx = 0
	}

	return c.nodes[c.sortedHashes[idx]]
}

func (c *ShardedClient) ChangeShardInfo(id string, si *ShardInfo) {
	c.Lock()
	defer c.Unlock()
	old := c.resources[id]
	c.resources[id] = si
	for k, v := range c.nodes {
		if v == old {
			c.nodes[k] = si
		}
	}
	old.client.Close()
}

func (c *ShardedClient) getKeyTag(key string) string {
	if c.tagPattern != nil {
		m := c.tagPattern.FindStringSubmatch(key)
		if len(m) > 1 {
			return m[1]
		}
	}
	return key
}

type ShardInfo struct {
	id     string
	name   string
	client Redis
	weight int
}

type Hashing interface {
	hash(key string) int64
}

type MurmurHash struct {
}

func (h *MurmurHash) hash(key string) int64 {
	return MurmurHashString(key)
}

type cmdsMap struct {
	mu sync.Mutex
	m  map[Redis][]redis.Cmder
}

func newCmdsMap() *cmdsMap {
	return &cmdsMap{
		m: make(map[Redis][]redis.Cmder),
	}
}

func (m *cmdsMap) Add(client Redis, cmds ...redis.Cmder) {
	m.mu.Lock()
	m.m[client] = append(m.m[client], cmds...)
	m.mu.Unlock()
}

func cmdsFirstErr(cmds []redis.Cmder) error {
	for _, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			return err
		}
	}
	return nil
}
