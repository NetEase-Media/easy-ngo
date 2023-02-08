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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
)

func BenchmarkHTTP_1Conn_1Delay(b *testing.B) {
	testRunBench(b, 1, 1)
}

func BenchmarkHTTP_1Conn_5Delay(b *testing.B) {
	testRunBench(b, 1, 5)
}

func BenchmarkHTTP_1Conn_50Delay(b *testing.B) {
	testRunBench(b, 1, 50)
}

func BenchmarkHTTP_5Conn_1Delay(b *testing.B) {
	testRunBench(b, 5, 1)
}

func BenchmarkHTTP_5Conn_5Delay(b *testing.B) {
	testRunBench(b, 5, 5)
}

func BenchmarkHTTP_5Conn_50Delay(b *testing.B) {
	testRunBench(b, 5, 50)
}

func BenchmarkHTTP_100Conn_1Delay(b *testing.B) {
	testRunBench(b, 100, 1)
}

func BenchmarkHTTP_100Conn_5Delay(b *testing.B) {
	testRunBench(b, 100, 5)
}

func BenchmarkHTTP_100Conn_50Delay(b *testing.B) {
	testRunBench(b, 100, 50)
}

func testRunBench(b *testing.B, n int, blockMS int) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * time.Duration(blockMS))
	}))
	defer s.Close()

	c, _ := newWithOption(&Option{
		MaxConnsPerHost:    n,
		MaxConnWaitTimeout: time.Second * 10,
	}, &xfmt.XFmt{}, nil, nil)
	b.SetParallelism(20)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := c.Get(s.URL).doInternal()
			CheckError(err)
		}
	})
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
