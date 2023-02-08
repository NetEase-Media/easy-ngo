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

/* 需要重构
var (
	crdb *ClusterClient
)

func init() {
	crdb = NewClusterClient(&ClusterOptions{
		ClusterOptions: redis.ClusterOptions{
			Addr: []string{"127.0.0.1:30021""},
		},
		Name: "test",
	})
	status, err := crdb.Ping(context.Background()).Result()
	fmt.Println(status, err)
}

func TestClusterClient_Set(t *testing.T) {
	k1 := generateKey()
	str, _ := crdb.SetEX(context.Background(), k1, "1000", time.Second*20).Result()
	if str != "OK" {
		t.Fatal("SetEX not valid", str)
	}
	str2, _ := crdb.Get(context.Background(), k1).Result()
	if str2 != "1000" {
		t.Fatal("Get not valid", str2)
	}
}

func TestClusterClient_Del(t *testing.T) {
	k1 := generateKey()
	crdb.SetEX(context.Background(), k1, "1000", time.Second*60).Result()
	delNum, err := crdb.Del(context.Background(), k1).Result()
	if delNum != 1 {
		t.Fatal("Del not valid", delNum, err)
	}
}

//TODO Exists 多个key会报错？
func TestClusterClient_Exists(t *testing.T) {
	k1 := generateKey()

	if num, err := crdb.Exists(context.Background(), k1).Result(); num != 0 {
		t.Fatal("got not valid", err)
	}
	crdb.SetEX(context.Background(), k1, "1", time.Second*60).Result()
	if num, err := crdb.Exists(context.Background(), k1).Result(); num != 1 {
		t.Fatal("got not valid", err)
	}
}

func TestClusterClient_ZRange(t *testing.T) {
	k1 := generateKey()
	crdb.ZAdd(context.Background(), k1, &redis.Z{Score: 1.55, Member: "key1"}).Result()
	crdb.ZAdd(context.Background(), k1, &redis.Z{Score: 1.56, Member: "key2"}).Result()
	crdb.Expire(context.Background(), k1, time.Minute)

	vals, err := crdb.ZRange(context.Background(), k1, 0, -1).Result()
	if len(vals) != 2 {
		t.Fatal("got not valid", err)
	}
	if delNum, _ := crdb.Del(context.Background(), k1).Result(); delNum < 1 {
		t.Fatal("del not valid", err)
	}
}
*/
