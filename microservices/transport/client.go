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

package transport

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/microservices/internal/xrpc"
	"github.com/NetEase-Media/easy-ngoservices/middleware"
)

var (
	clientInfoKey struct{}
)

type CallOption func(o *options)

// Client is a transport client interface.
type Client interface {
	Use(selector string, mws ...middleware.Middleware)
	//Invoke(ctx context.Context, method string, req interface{}, resp interface{}, opts ...CallOption) error
	Close() error
}

type options struct {
}

type ClientInfo struct {
	Type   string
	Target string
	Op     *xrpc.Operation
}

func NewClientContext(ctx context.Context, info ClientInfo) context.Context {
	return context.WithValue(ctx, clientInfoKey, info)
}

func FromClientContext(ctx context.Context) (info ClientInfo, ok bool) {
	info, ok = ctx.Value(clientInfoKey).(ClientInfo)
	return
}
