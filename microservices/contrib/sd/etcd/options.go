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
	"time"

	"github.com/NetEase-Media/easy-ngo/xlog"
)

// Option is a function that configures the service discovery.
type Option func(o *options)

// WithNamespace returns a new Option that sets the namespace.
func WithNamespace(namespace string) Option {
	return func(o *options) {
		o.namespace = namespace
	}
}

// WithTTL returns a new Option that sets the ttl.
func WithTTL(ttl time.Duration) Option {
	return func(o *options) {
		if ttl > 0 {
			o.ttl = ttl
		}
	}
}

// WithLogger returns a new Option that sets the logger.
func WithLogger(logger xlog.Logger) Option {
	return func(o *options) {
		o.log = logger
	}
}

// defaultOptions returns the default options.
func defaultOptions() *options {
	return &options{
		namespace: "microservices",
		ttl:       time.Second * 15,
		log:       xlog.NewNopLogger(),
	}
}

// options is the options for service discovery.
type options struct {
	namespace string
	ttl       time.Duration
	log       xlog.Logger
}
