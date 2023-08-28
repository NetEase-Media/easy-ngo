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

package xfasthttp

import (
	"github.com/valyala/fasthttp"
)

const (
	// 最大重定向次数
	defaultMaxRedirectsCount = 16
)

type Xfasthttp struct {
	client *fasthttp.Client
}

func New(c *Config) (*Xfasthttp, error) {
	client := &fasthttp.Client{
		Name:                      c.Name,
		NoDefaultUserAgentHeader:  c.NoDefaultUserAgentHeader,
		TLSConfig:                 c.TLSConfig, // TODO: TLS系列配置需要另外分离
		MaxConnsPerHost:           c.MaxConnsPerHost,
		MaxIdleConnDuration:       c.MaxIdleConnDuration,
		MaxConnDuration:           c.MaxConnDuration,
		MaxIdemponentCallAttempts: c.MaxIdemponentCallAttempts,
		ReadBufferSize:            c.ReadBufferSize,
		WriteBufferSize:           c.WriteBufferSize,
		ReadTimeout:               c.ReadTimeout,
		WriteTimeout:              c.WriteTimeout,
		MaxResponseBodySize:       c.MaxResponseBodySize,
		MaxConnWaitTimeout:        c.MaxConnWaitTimeout,
	}
	fhttp := &Xfasthttp{
		client: client,
	}
	return fhttp, nil
}

// Get 调用http客户端的GET方法
func (f *Xfasthttp) Get(url string) *DataFlow {
	df := f.newDataFlow()
	return df.newMethod(fasthttp.MethodGet, url)
}

// Post 调用http客户端的POST方法
func (f *Xfasthttp) Post(url string) *DataFlow {
	df := f.newDataFlow()
	return df.newMethod(fasthttp.MethodPost, url)
}

// Put 调用http客户端的PUT方法
func (f *Xfasthttp) Put(url string) *DataFlow {
	df := f.newDataFlow()
	return df.newMethod(fasthttp.MethodPut, url)
}

// Delete 调用http客户端的DELETE方法
func (f *Xfasthttp) Delete(url string) *DataFlow {
	df := f.newDataFlow()
	return df.newMethod(fasthttp.MethodDelete, url)
}

// Patch 调用http客户端的PATCH方法
func (f *Xfasthttp) Patch(url string) *DataFlow {
	df := f.newDataFlow()
	return df.newMethod(fasthttp.MethodPatch, url)
}

func (f *Xfasthttp) Close() {
	f.client.CloseIdleConnections()
}

func (f *Xfasthttp) newDataFlow() *DataFlow {
	df := newDataFlow(f)
	return df
}
