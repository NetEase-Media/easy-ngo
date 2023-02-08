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

package sd

import "context"

// Registrar registers instance information to a service discovery system
type Registrar interface {
	Register(ctx context.Context, service *ServiceInfo) error
	Deregister(ctx context.Context, service *ServiceInfo) error
}

// Discovery discovers instance information from a service discovery system
type Discovery interface {
	GetService(ctx context.Context, serviceName string) ([]*ServiceInfo, error)
	NewWatcher(ctx context.Context, serviceName string) (Watcher, error)
}

type Watcher <-chan []*Update
type Operation uint8

const (
	Add Operation = iota
	Delete
)

type Update struct {
	Op          Operation
	Key         string
	ServiceInfo *ServiceInfo
}

type ServiceInfo struct {
	Name     string
	Scheme   string
	Addr     string
	Metadata map[string]string
}
