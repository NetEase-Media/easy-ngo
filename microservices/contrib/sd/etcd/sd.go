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

package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/NetEase-Media/easy-ngo/microservices/sd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var _ sd.Registrar = (*ServiceDiscovery)(nil)
var _ sd.Discovery = (*ServiceDiscovery)(nil)

// New creates a new etcd service discovery.
func New(ctx context.Context, client *clientv3.Client, optFns ...Option) *ServiceDiscovery {
	sd := ServiceDiscovery{
		ctx:    ctx,
		client: client,
		opts:   defaultOptions(),
		kvs:    sync.Map{},
		lmu:    &sync.RWMutex{},
	}
	for _, o := range optFns {
		o(sd.opts)
	}
	return &sd
}

// ServiceDiscovery is a etcd service discovery.
type ServiceDiscovery struct {
	ctx    context.Context
	client *clientv3.Client

	opts *options

	kvs     sync.Map
	lmu     *sync.RWMutex
	leaseID clientv3.LeaseID
	once    sync.Once
}

// Register registers a service.
func (s *ServiceDiscovery) Register(ctx context.Context, service *sd.ServiceInfo) error {
	var key, value string
	key = s.getKey(service.Name, service.Addr)
	if v, err := json.Marshal(service); err != nil {
		return err
	} else {
		value = string(v)
	}

	leaseID, err := s.getOrGrantLeaseID(ctx)
	if err != nil {
		return err
	}
	_, err = s.client.Put(ctx, key, value, clientv3.WithLease(leaseID))
	if err != nil {
		return err
	}

	s.kvs.Store(key, value)
	s.once.Do(func() {
		go s.keepalive(s.ctx)
	})
	return nil
}

// Deregister deregisters a service.
func (s *ServiceDiscovery) Deregister(ctx context.Context, service *sd.ServiceInfo) error {
	key := s.getKey(service.Name, service.Addr)
	_, err := s.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	s.kvs.Delete(key)
	return nil
}

// GetService gets a service.
func (s *ServiceDiscovery) GetService(ctx context.Context, serviceName string) ([]*sd.ServiceInfo, error) {
	resp, err := s.client.Get(ctx, s.getPrefix(serviceName), clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}

	sis := make([]*sd.ServiceInfo, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var si sd.ServiceInfo
		if err := json.Unmarshal(kv.Value, &si); err != nil {
			continue
		}
		sis = append(sis, &si)
	}
	return sis, nil
}

// NewWatcher creates a new watcher.
func (s *ServiceDiscovery) NewWatcher(ctx context.Context, serviceName string) (sd.Watcher, error) {
	resp, err := s.client.Get(ctx, s.getPrefix(serviceName), clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}

	ups := make([]*sd.Update, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var si sd.ServiceInfo
		if err := json.Unmarshal(kv.Value, &si); err != nil {
			continue
		}
		ups = append(ups, &sd.Update{Op: sd.Add, Key: s.getKey(si.Name, si.Addr), ServiceInfo: &si})
	}

	upch := make(chan []*sd.Update, 1)
	if len(ups) > 0 {
		upch <- ups
	}
	go s.watch(ctx, serviceName, resp.Header.Revision+1, upch)
	return upch, nil
}

// getOrGrantLeaseID gets or grants a lease ID.
func (s *ServiceDiscovery) getOrGrantLeaseID(ctx context.Context) (clientv3.LeaseID, error) {
	s.lmu.Lock()
	defer s.lmu.Unlock()

	if s.leaseID != 0 {
		return s.leaseID, nil
	}

	grant, err := s.client.Grant(ctx, int64(s.opts.ttl.Seconds()))
	if err != nil {
		return 0, err
	}

	s.leaseID = grant.ID
	return grant.ID, nil
}

// getLeaseID gets a lease ID.
func (s *ServiceDiscovery) getLeaseID() clientv3.LeaseID {
	s.lmu.RLock()
	defer s.lmu.RUnlock()
	return s.leaseID
}

// setLeaseID sets a lease ID.
func (s *ServiceDiscovery) setLeaseID(leaseID clientv3.LeaseID) {
	s.lmu.Lock()
	defer s.lmu.Unlock()
	s.leaseID = leaseID
}

// keepalive keeps alive.
func (s *ServiceDiscovery) keepalive(ctx context.Context) {
	kac, err := s.client.KeepAlive(ctx, s.getLeaseID())
	if err != nil {
		s.setLeaseID(0)
	}

	for {
		if ctx.Err() != nil {
			return
		}
		if s.getLeaseID() == 0 {
			done := make(chan struct{}, 1)
			cancelCtx, cancel := context.WithCancel(ctx)
			go func() {
				defer func() {
					done <- struct{}{}
				}()
				leaseID, err := s.getOrGrantLeaseID(cancelCtx)
				if err != nil {
					s.opts.log.Errorf("get or grant LeaseID failed err=%v", err)
					return
				}
				ops := make([]clientv3.Op, 0, 3)
				s.kvs.Range(func(key, value interface{}) bool {
					ops = append(ops, clientv3.OpPut(key.(string), value.(string), clientv3.WithLease(leaseID)))
					return true
				})
				_, err = s.client.KV.Txn(ctx).Then(ops...).Commit()
				if err != nil {
					s.opts.log.Errorf("register kv failed err=%v", err)
				}
			}()

			select {
			case <-time.After(time.Second * 5):
				continue
			case <-done:
				if s.getLeaseID() == 0 {
					continue
				}
			}

			cancel()

			kac, err = s.client.KeepAlive(ctx, s.getLeaseID())
			if err != nil {
				s.setLeaseID(0)
				time.Sleep(time.Second * 3)
				continue
			}
		}

		select {
		case _, ok := <-kac:
			if !ok {
				s.setLeaseID(0)
			}
		case <-ctx.Done():
			return
		}
	}
}

// watch watches a service.
func (s *ServiceDiscovery) watch(ctx context.Context, serviceName string, rev int64, upch chan []*sd.Update) {
	defer close(upch)

	keyPrefix := s.getPrefix(serviceName)
	opts := []clientv3.OpOption{clientv3.WithRev(rev), clientv3.WithPrefix()}
	wch := s.client.Watch(ctx, keyPrefix, opts...)
	for {
		select {
		case <-ctx.Done():
			return
		case wresp, ok := <-wch:
			if !ok {
				s.opts.log.Warnf("watch closed keyPrefix=%s", keyPrefix)
				return
			}
			if wresp.Err() != nil {
				s.opts.log.Errorf("watch failed keyPrefix=%s err=%s", keyPrefix, wresp.Err())
				return
			}

			deltaUps := make([]*sd.Update, 0, len(wresp.Events))
			for _, e := range wresp.Events {
				var si sd.ServiceInfo
				var err error
				var op sd.Operation
				switch e.Type {
				case clientv3.EventTypePut:
					err = json.Unmarshal(e.Kv.Value, &si)
					op = sd.Add
					if err != nil {
						s.opts.log.Warnf("unmarshal service info failed key=%s err=%v", string(e.Kv.Key), err)
						continue
					}
				case clientv3.EventTypeDelete:
					op = sd.Delete
				default:
					continue
				}
				up := &sd.Update{Op: op, Key: string(e.Kv.Key), ServiceInfo: &si}
				deltaUps = append(deltaUps, up)
			}
			if len(deltaUps) > 0 {
				upch <- deltaUps
			}
		}
	}
}

// getPrefix gets a prefix.
func (s *ServiceDiscovery) getPrefix(serviceName string) string {
	return fmt.Sprintf("%s/%s/", s.opts.namespace, serviceName)
}

// getKey gets a key.
func (s *ServiceDiscovery) getKey(serviceName, addr string) string {
	return fmt.Sprintf("%s/%s/%s", s.opts.namespace, serviceName, addr)
}
