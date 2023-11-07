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
	"net/http"
)

type METHOD string

const (
	GET     METHOD = http.MethodGet
	HEAD           = http.MethodHead
	POST           = http.MethodPost
	PUT            = http.MethodPut
	PATCH          = http.MethodPatch
	DELETE         = http.MethodDelete
	CONNECT        = http.MethodConnect
	OPTIONS        = http.MethodOptions
	TRACE          = http.MethodTrace
)

type Server interface {
	Serve() error
	Shutdown() error
	Healthz() bool
	Init() error

	GET(relativePath string, handler any)
	POST(relativePath string, handler any)
	PUT(relativePath string, handler any)
	DELETE(relativePath string, handler any)
	PATCH(relativePath string, handler any)
	HEAD(relativePath string, handler any)
	OPTIONS(relativePath string, handler any)
}
