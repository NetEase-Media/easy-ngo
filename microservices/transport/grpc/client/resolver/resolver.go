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

package resolver

import (
	"context"
	"strings"
	"sync"

	"github.com/NetEase-Media/easy-ngo/microservices/sd"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/codes"
	gresolver "google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

var _ gresolver.Builder = (*builder)(nil)
var _ gresolver.Resolver = (*resolver)(nil)

// NewBuilder creates a resolver builder.
func NewBuilder(d sd.Discovery) gresolver.Builder {
	return &builder{discovery: d}
}

type builder struct {
	discovery sd.Discovery
}

func (b builder) Build(target gresolver.Target, cc gresolver.ClientConn, opts gresolver.BuildOptions) (gresolver.Resolver, error) {
	r := &resolver{
		target: target.URL.Path,
		cc:     cc,
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())

	watcher, err := b.discovery.NewWatcher(r.ctx, strings.TrimPrefix(r.target, "/"))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resolver: failed to new watcher: %s", err)
	}
	r.watcher = watcher

	r.wg.Add(1)
	go r.watch()
	return r, nil
}

func (b builder) Scheme() string {
	return "sd"
}

type resolver struct {
	target  string
	cc      gresolver.ClientConn
	watcher sd.Watcher
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func (r *resolver) watch() {
	defer r.wg.Done()

	allUps := make(map[string]*sd.Update)
	for {
		select {
		case <-r.ctx.Done():
			return
		case ups, ok := <-r.watcher:
			if !ok {
				return
			}

			for _, up := range ups {
				switch up.Op {
				case sd.Add:
					allUps[up.Key] = up
				case sd.Delete:
					delete(allUps, up.Key)
				}
			}

			addrs := convertToGRPCAddress(allUps)
			r.cc.UpdateState(gresolver.State{Addresses: addrs})
		}
	}
}

func convertToGRPCAddress(ups map[string]*sd.Update) []gresolver.Address {
	var addrs []gresolver.Address
	for _, up := range ups {
		addr := gresolver.Address{
			Addr:       up.ServiceInfo.Addr,
			Attributes: parseAttributes(up.ServiceInfo.Metadata),
		}
		addrs = append(addrs, addr)
	}
	return addrs
}

// ResolveNow is a no-op here.
// It's just a hint, resolver can ignore this if it's not necessary.
func (r *resolver) ResolveNow(gresolver.ResolveNowOptions) {}

func (r *resolver) Close() {
	r.cancel()
	r.wg.Wait()
}

func parseAttributes(md map[string]string) *attributes.Attributes {
	var a *attributes.Attributes
	for k, v := range md {
		if a == nil {
			a = attributes.New(k, v)
		} else {
			a = a.WithValue(k, v)
		}
	}
	return a
}
