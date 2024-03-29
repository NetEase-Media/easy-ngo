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

package app

import (
	"context"
	"sync"
)

var (
	globalPlugins = make(map[Status][]func(ctx context.Context) error)
	mu            = sync.RWMutex{}
)

func RegisterPlugin(status Status, fns ...func(ctx context.Context) error) {
	mu.Lock()
	defer mu.Unlock()
	globalPlugins[status] = append(globalPlugins[status], fns...)
}

func GetFns(status Status) []func(ctx context.Context) error {
	return globalPlugins[status]
}
