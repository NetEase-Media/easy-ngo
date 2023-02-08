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

package httplib

import (
	"github.com/NetEase-Media/easy-ngo/observability/metrics"
	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/valyala/fasthttp"
)

const (
	// 最大重定向次数
	defaultMaxRedirectsCount = 16
)

type HttpClient struct {
	client  *fasthttp.Client
	opt     Option
	logger  xlog.Logger
	metrics metrics.Provider
	tracer  tracer.Provider
}

func New(opt *Option, logger xlog.Logger, metrics metrics.Provider, tracer tracer.Provider) (*HttpClient, error) {
	return newWithOption(opt, logger, metrics, tracer)
}

func newWithOption(opt *Option, logger xlog.Logger, metrics metrics.Provider, tracer tracer.Provider) (*HttpClient, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}
	client := &fasthttp.Client{
		Name:                      opt.Name,
		NoDefaultUserAgentHeader:  opt.NoDefaultUserAgentHeader,
		TLSConfig:                 opt.TLSConfig, // TODO: TLS系列配置需要另外分离
		MaxConnsPerHost:           opt.MaxConnsPerHost,
		MaxIdleConnDuration:       opt.MaxIdleConnDuration,
		MaxConnDuration:           opt.MaxConnDuration,
		MaxIdemponentCallAttempts: opt.MaxIdemponentCallAttempts,
		ReadBufferSize:            opt.ReadBufferSize,
		WriteBufferSize:           opt.WriteBufferSize,
		ReadTimeout:               opt.ReadTimeout,
		WriteTimeout:              opt.WriteTimeout,
		MaxResponseBodySize:       opt.MaxResponseBodySize,
		MaxConnWaitTimeout:        opt.MaxConnWaitTimeout,
	}
	c := &HttpClient{
		client:  client,
		opt:     *opt,
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
	c.initMetrics()
	return c, nil
}

// Get 调用http客户端的GET方法
func (c *HttpClient) Get(url string) *DataFlow {
	df := c.newDataFlow()
	return df.newMethod(fasthttp.MethodGet, url)
}

// Post 调用http客户端的POST方法
func (c *HttpClient) Post(url string) *DataFlow {
	df := c.newDataFlow()
	return df.newMethod(fasthttp.MethodPost, url)
}

// Put 调用http客户端的PUT方法
func (c *HttpClient) Put(url string) *DataFlow {
	df := c.newDataFlow()
	return df.newMethod(fasthttp.MethodPut, url)
}

// Delete 调用http客户端的DELETE方法
func (c *HttpClient) Delete(url string) *DataFlow {
	df := c.newDataFlow()
	return df.newMethod(fasthttp.MethodDelete, url)
}

// Patch 调用http客户端的PATCH方法
func (c *HttpClient) Patch(url string) *DataFlow {
	df := c.newDataFlow()
	return df.newMethod(fasthttp.MethodPatch, url)
}

func (c *HttpClient) Close() {
	c.client.CloseIdleConnections()
}

func (c *HttpClient) newDataFlow() *DataFlow {
	df := newDataFlow(c, c.logger, c.metrics, c.tracer)
	if c.opt.EnableTracer {
		df.WrapDoFunc(tracingDoFuc)
	}
	return df
}
