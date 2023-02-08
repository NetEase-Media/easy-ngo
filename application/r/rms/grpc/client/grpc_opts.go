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

package client

import (
	"github.com/NetEase-Media/easy-ngo/application/r/rms/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func getGRPCOptions(opts *api.GRPCClientOptions) []grpc.DialOption {
	options := make([]grpc.DialOption, 0, 8)
	if opts.InitialWindowSize > 0 {
		options = append(options, grpc.WithInitialWindowSize(opts.InitialWindowSize))
	}
	if opts.InitialConnWindowSize > 0 {
		options = append(options, grpc.WithInitialConnWindowSize(opts.InitialConnWindowSize))
	}
	if opts.MaxHeaderListSize > 0 {
		options = append(options, grpc.WithMaxHeaderListSize(opts.MaxHeaderListSize))
	}
	if opts.ReadBufferSize > 0 {
		options = append(options, grpc.WithReadBufferSize(int(opts.ReadBufferSize)))
	}
	if opts.WriteBufferSize > 0 {
		options = append(options, grpc.WithWriteBufferSize(int(opts.WriteBufferSize)))
	}
	if opts.Authority != "" {
		options = append(options, grpc.WithAuthority(opts.Authority))
	}
	if opts.Block {
		options = append(options, grpc.WithBlock())
	}
	if opts.UserAgent != "" {
		options = append(options, grpc.WithUserAgent(opts.UserAgent))
	}
	if opts.DisableRetry {
		options = append(options, grpc.WithDisableRetry())
	}
	if opts.ConnectParams != nil {
		var connectParams grpc.ConnectParams
		if opts.ConnectParams.MinConnectTimeout != nil && opts.ConnectParams.MinConnectTimeout.AsDuration() > 0 {
			connectParams.MinConnectTimeout = opts.ConnectParams.MinConnectTimeout.AsDuration()
		}
		if opts.ConnectParams.BaseDelay != nil {
			connectParams.Backoff.BaseDelay = opts.ConnectParams.BaseDelay.AsDuration()
		}
		if opts.ConnectParams.MaxDelay != nil {
			connectParams.Backoff.MaxDelay = opts.ConnectParams.MaxDelay.AsDuration()
		}
		if opts.ConnectParams.Jitter > 0 {
			connectParams.Backoff.Jitter = float64(opts.ConnectParams.Jitter)
		}
		if opts.ConnectParams.Multiplier > 0 {
			connectParams.Backoff.Multiplier = float64(opts.ConnectParams.Multiplier)
		}
		options = append(options, grpc.WithConnectParams(connectParams))
	}
	if opts.KeepaliveParams != nil {
		var keepaliveParams keepalive.ClientParameters
		if opts.KeepaliveParams.Time != nil {
			keepaliveParams.Time = opts.KeepaliveParams.Time.AsDuration()
		}
		if opts.KeepaliveParams.Timeout != nil {
			keepaliveParams.Timeout = opts.KeepaliveParams.Timeout.AsDuration()
		}
		if opts.KeepaliveParams.PermitWithoutStream {
			keepaliveParams.PermitWithoutStream = true
		}
		options = append(options, grpc.WithKeepaliveParams(keepaliveParams))
	}

	return options
}
