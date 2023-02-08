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

package hooks

import (
	"context"
	"sync"
)

type Stage int

const (
	Initialize Stage = iota
	Start
	Stop
	Online
	Offline
	HealthCheck
)

var (
	globalHooks = make(map[Stage][]func(ctx context.Context) error)
	mu          = sync.RWMutex{}
)

func Register(stage Stage, fns ...func(ctx context.Context) error) {
	mu.Lock()
	defer mu.Unlock()
	globalHooks[stage] = append(globalHooks[stage], fns...)
}

func GetFns(stage Stage) []func(ctx context.Context) error {
	return globalHooks[stage]
}
