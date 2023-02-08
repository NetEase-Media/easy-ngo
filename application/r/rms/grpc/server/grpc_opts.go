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

package server

import (
	"github.com/NetEase-Media/easy-ngo/application/r/rms/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func getGRPCOptions(opts *api.GRPCServerOptions) []grpc.ServerOption {
	options := make([]grpc.ServerOption, 0, 8)
	if opts.InitialConnWindowSize > 0 {
		options = append(options, grpc.InitialConnWindowSize(opts.InitialConnWindowSize))
	}
	if opts.InitialWindowSize > 0 {
		options = append(options, grpc.InitialWindowSize(opts.InitialWindowSize))
	}
	if opts.MaxConcurrentStreams > 0 {
		options = append(options, grpc.MaxConcurrentStreams(opts.MaxConcurrentStreams))
	}
	if opts.MaxHeaderListSize > 0 {
		options = append(options, grpc.MaxHeaderListSize(opts.MaxHeaderListSize))
	}
	if opts.MaxReceiveMessageSize > 0 {
		options = append(options, grpc.MaxRecvMsgSize(int(opts.MaxReceiveMessageSize)))
	}
	if opts.MaxSendMessageSize > 0 {
		options = append(options, grpc.MaxSendMsgSize(int(opts.MaxSendMessageSize)))
	}
	if opts.ReadBufferSize > 0 {
		options = append(options, grpc.ReadBufferSize(int(opts.ReadBufferSize)))
	}
	if opts.WriteBufferSize > 0 {
		options = append(options, grpc.WriteBufferSize(int(opts.WriteBufferSize)))
	}
	if opts.ConnectionTimeout != nil && opts.ConnectionTimeout.AsDuration() > 0 {
		options = append(options, grpc.ConnectionTimeout(opts.ConnectionTimeout.AsDuration()))
	}
	if opts.KeepalivePolicy != nil {
		var keepalivePolicy keepalive.EnforcementPolicy
		if opts.KeepalivePolicy.MinTime != nil && opts.KeepalivePolicy.MinTime.AsDuration() > 0 {
			keepalivePolicy.MinTime = opts.KeepalivePolicy.MinTime.AsDuration()
		}
		if opts.KeepalivePolicy.PermitWithoutStream {
			keepalivePolicy.PermitWithoutStream = opts.KeepalivePolicy.PermitWithoutStream
		}
		options = append(options, grpc.KeepaliveEnforcementPolicy(keepalivePolicy))
	}
	if opts.KeepaliveParams != nil {
		var parameters keepalive.ServerParameters
		if opts.KeepaliveParams.MaxConnectionIdle != nil && opts.KeepaliveParams.MaxConnectionIdle.AsDuration() > 0 {
			parameters.MaxConnectionIdle = opts.KeepaliveParams.MaxConnectionIdle.AsDuration()
		}
		if opts.KeepaliveParams.MaxConnectionAge != nil && opts.KeepaliveParams.MaxConnectionAge.AsDuration() > 0 {
			parameters.MaxConnectionAge = opts.KeepaliveParams.MaxConnectionAge.AsDuration()
		}
		if opts.KeepaliveParams.MaxConnectionAgeGrace != nil && opts.KeepaliveParams.MaxConnectionAgeGrace.AsDuration() > 0 {
			parameters.MaxConnectionAgeGrace = opts.KeepaliveParams.MaxConnectionAgeGrace.AsDuration()
		}
		if opts.KeepaliveParams.Time != nil && opts.KeepaliveParams.Timeout.AsDuration() > 0 {
			parameters.Time = opts.KeepaliveParams.Time.AsDuration()
		}
		if opts.KeepaliveParams.Timeout != nil && opts.KeepaliveParams.Timeout.AsDuration() > 0 {
			parameters.Timeout = opts.KeepaliveParams.Timeout.AsDuration()
		}
		options = append(options, grpc.KeepaliveParams(parameters))
	}
	return options
}
