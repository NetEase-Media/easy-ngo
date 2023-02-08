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

package timeout

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CallBackFunc func(*http.Request)

// Option for timeout
type Option func(*Options)

// WithTimeout set timeout
func WithTimeout(timeout time.Duration) Option {
	return func(t *Options) {
		if timeout > 0 {
			t.timeout = timeout
		}
	}
}

// WithHandler set handle process
func WithHandler(f gin.HandlerFunc) Option {
	return func(t *Options) {
		t.handler = f
	}
}

// Optional parameters
func WithErrorHttpCode(code int) Option {
	return func(t *Options) {
		t.errorHttpCode = code
	}
}

// Optional parameters
func WithDefaultMsg(s string) Option {
	return func(t *Options) {
		t.defaultMsg = s
	}
}

// Optional parameters
func WithCallBack(f CallBackFunc) Option {
	return func(t *Options) {
		t.callBack = f
	}
}

// Optional parameters
func WithErrorHandler(f gin.HandlerFunc) Option {
	return func(t *Options) {
		t.errorHandler = f
	}
}

// Options struct
type Options struct {
	timeout       time.Duration
	handler       gin.HandlerFunc
	errorHttpCode int
	defaultMsg    string
	errorHandler  gin.HandlerFunc
	callBack      CallBackFunc
}
