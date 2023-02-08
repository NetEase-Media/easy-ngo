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

import "github.com/NetEase-Media/easy-ngo/xlog/xfmt"

var defaultHttpClient = Default()

func Default() *HttpClient {
	opt := DefaultOption()
	opt.EnableTracer = true
	client, _ := newWithOption(opt, &xfmt.XFmt{}, nil, nil)
	return client
}

func DefaultHttpClient() *HttpClient {
	return defaultHttpClient
}

func SetDefaultHttpClient(client *HttpClient) {
	defaultHttpClient = client
}

// Get 调用默认http客户端的GET方法
func Get(url string) *DataFlow {
	return defaultHttpClient.Get(url)
}

// Post 调用默认http客户端的POST方法
func Post(url string) *DataFlow {
	return defaultHttpClient.Post(url)
}

// Put 调用默认http客户端的PUT方法
func Put(url string) *DataFlow {
	return defaultHttpClient.Put(url)
}

// Delete 调用默认http客户端的DELETE方法
func Delete(url string) *DataFlow {
	return defaultHttpClient.Delete(url)
}

// Patch 调用默认http客户端的PATCH方法
func Patch(url string) *DataFlow {
	return defaultHttpClient.Patch(url)
}

func Close() {
	defaultHttpClient.Close()
}
