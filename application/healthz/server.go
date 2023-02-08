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

package healthz

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngo/application/util"
	conf "github.com/NetEase-Media/easy-ngo/config"
)

type HealthzServer struct {
	active      bool
	activeMutex sync.Mutex

	mux *http.ServeMux

	Port            int
	OnlinePath      string
	OfflinePath     string
	HealthCheckPath string
}

func New() *HealthzServer {
	server := Default()
	conf.Get("ngo.app.healthz", server)
	return server
}

func Default() *HealthzServer {
	return &HealthzServer{
		Port:            18888,
		OnlinePath:      "/healthz/online",
		OfflinePath:     "/healthz/offline",
		HealthCheckPath: "/healthz/check",
	}
}

func (server *HealthzServer) init() {
	server.mux = http.NewServeMux()
	server.mux.HandleFunc(server.OnlinePath, server.OnlineHandler)
	server.mux.HandleFunc(server.OfflinePath, server.OfflineHandler)
	server.mux.HandleFunc(server.HealthCheckPath, server.HealthCheckHandler)
}

func (server *HealthzServer) Serve() error {
	server.init()
	http.ListenAndServe(fmt.Sprintf(":%d", server.Port), server.mux)
	return nil
}

func (server *HealthzServer) OnlineHandler(w http.ResponseWriter, r *http.Request) {
	server.activeMutex.Lock()
	defer server.activeMutex.Unlock()
	ctx := context.Background()
	fs := hooks.GetFns(hooks.Online)
	for i := range fs {
		if err := fs[i](ctx); err != nil {
			util.CheckError(err)
		}
	}
	server.active = true
}

func (server *HealthzServer) OfflineHandler(w http.ResponseWriter, r *http.Request) {
	server.activeMutex.Lock()
	defer server.activeMutex.Unlock()
	ctx := context.Background()
	fs := hooks.GetFns(hooks.Offline)
	for i := range fs {
		if err := fs[i](ctx); err != nil {
			util.CheckError(err)
		}
	}
	server.active = false
}

func (server *HealthzServer) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	server.activeMutex.Lock()
	active := server.active
	server.activeMutex.Unlock()
	if !active {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
