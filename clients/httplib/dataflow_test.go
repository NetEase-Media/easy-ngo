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
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"testing"
	"time"

	"github.com/NetEase-Media/easy-ngo/clients/xsentinel"
	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func newTestHttpClient() *HttpClient {
	c := &HttpClient{
		client: &fasthttp.Client{},
		opt:    *DefaultOption(),
	}
	return c
}

func testNewDataFlow() *DataFlow {
	lg, _ := xfmt.Default()
	df := newDataFlow(newTestHttpClient(), lg, nil, nil)

	return df
}

type testResponse struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func TestDataFlowGet(t *testing.T) {
	cases := []struct {
		url        string
		requestUrl string
	}{
		{
			url:        "http://www.google.com",
			requestUrl: "http://www.google.com",
		},
		{
			url:        "https://www.google.com",
			requestUrl: "https://www.google.com",
		},
	}

	df := testNewDataFlow()
	for _, c := range cases {
		df.newMethod(fasthttp.MethodGet, c.url)
		assert.Equal(t, string(df.req.Header.Method()), fasthttp.MethodGet)
		assert.Equal(t, c.requestUrl, string(df.req.RequestURI()))
	}
}

func TestDataFlowSetQuery(t *testing.T) {
	q := Query{
		"a": []string{"1"},
		"b": []string{"2"},
		"c": []string{"3", "4"},
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "a=1&b=2&c=3&c=4", r.URL.Query().Encode())
	}))
	df := testNewDataFlow()
	_, err := df.newMethod(fasthttp.MethodGet, s.URL).SetQuery(q).Do(context.Background())
	assert.Nil(t, err)
}

func TestDataFlowAddQuery(t *testing.T) {
	q := Query{
		"a": []string{"1"},
		"b": []string{"2"},
		"c": []string{"3", "4"},
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "a=1&b=2&c=3&c=4", r.URL.Query().Encode())
	}))
	df := testNewDataFlow().newMethod(fasthttp.MethodGet, s.URL)
	for k, vs := range q {
		for _, v := range vs {
			df.AddQuery(k, v)
		}
	}
	_, err := df.Do(context.Background())
	assert.Nil(t, err)
}

type testJsonBody struct {
	A string  `json:"a"`
	B int     `json:"b"`
	C float64 `json:"c"`
}

func TestDataFlowBindJson(t *testing.T) {

	mustMashal := func(obj interface{}) []byte {
		b, err := json.Marshal(obj)
		assert.Nil(t, err)
		return b
	}

	cases := []struct {
		obj         interface{}
		body        []byte
		expectedObj interface{}
		hasError    bool
	}{
		{
			obj: &testJsonBody{},
			body: mustMashal(&testJsonBody{
				A: "dddd",
				B: 50,
				C: 60,
			}),
			expectedObj: &testJsonBody{
				A: "dddd",
				B: 50,
				C: 60,
			},
		},
		{
			obj:         new(int),
			body:        []byte("10"),
			expectedObj: testNewInt(10),
		},
		{
			obj: new(int),
			body: mustMashal(&testJsonBody{
				A: "int",
				B: 20,
				C: 10,
			}),
			hasError: true,
		},
	}

	for _, c := range cases {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(205)
			_, err := w.Write(c.body)
			assert.Nil(t, err)
		}))

		df := testNewDataFlow()
		statusCode, err := df.newMethod(fasthttp.MethodGet, s.URL).BindJson(&c.obj).Do(context.Background())
		if c.hasError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.EqualValues(t, c.expectedObj, c.obj)
			assert.Equal(t, 205, statusCode)
		}

		s.Close()
	}
}

func TestDataFlowBindInt(t *testing.T) {
	cases := []struct {
		body        []byte
		expectedObj int
	}{
		{
			body:        []byte("10"),
			expectedObj: 10,
		},
		{
			body:        []byte("20"),
			expectedObj: 20,
		},
	}

	for _, c := range cases {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(206)
			_, err := w.Write(c.body)
			assert.Nil(t, err)
		}))

		var i int
		df := testNewDataFlow()
		statusCode, err := df.newMethod(fasthttp.MethodGet, s.URL).BindInt(&i).Do(context.Background())
		assert.Nil(t, err)
		assert.EqualValues(t, c.expectedObj, i)
		assert.Equal(t, 206, statusCode)
		s.Close()
	}
}

func TestDataFlowAddHeader(t *testing.T) {
	header := H{
		"k1": {"v1", "v2", "v6"},
		"k3": {"v3", "v8", "v10"},
		"k7": {"v3", "v5", "v7"},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range header {
			values := r.Header.Values(k)
			sort.Strings(v)
			sort.Strings(values)
			assert.EqualValues(t, v, values)
		}
	}))
	defer s.Close()

	df := testNewDataFlow()
	_, err := df.newMethod(fasthttp.MethodGet, s.URL).AddHeader(header).Do(context.Background())
	assert.Nil(t, err)
}

func TestDataFlowAddHeaderKV(t *testing.T) {
	header := H{
		"K1": {"v1", "v2", "v6"},
		"K3": {"v3", "v8", "v10"},
		"K7": {"v3", "v5", "v7"},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range header {
			values := r.Header.Values(k)
			sort.Strings(v)
			sort.Strings(values)
			assert.EqualValues(t, v, values)
		}
	}))
	defer s.Close()

	df := testNewDataFlow()
	df = df.newMethod(fasthttp.MethodGet, s.URL)
	for key, values := range header {
		df = df.AddHeaderKV(key, values...)
	}
	_, err := df.Do(context.Background())
	assert.Nil(t, err)
}

func TestDataFlowBindHeader(t *testing.T) {
	header := H{
		"K1": {"v1", "v2", "v6"},
		"K3": {"v3", "v8", "v10"},
		"K7": {"v3", "v5", "v7"},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, arr := range header {
			for _, v := range arr {
				w.Header().Add(k, v)
			}
		}
	}))
	defer s.Close()

	bindHead := make(H)
	df := testNewDataFlow()
	_, err := df.newMethod(fasthttp.MethodGet, s.URL).BindHeader(bindHead).Do(context.Background())
	assert.Nil(t, err)

	for key, values := range header {
		sort.Strings(values)
		bindValues := bindHead[key]
		sort.Strings(bindValues)
		assert.EqualValues(t, bindValues, values)
	}

}

func TestDataFlowPost(t *testing.T) {
	cases := []struct {
		url        string
		requestUrl string
	}{
		{
			url:        "http://www.google.com",
			requestUrl: "http://www.google.com",
		},
		{
			url:        "https://www.google.com",
			requestUrl: "https://www.google.com",
		},
	}

	df := testNewDataFlow()
	for _, c := range cases {
		df.newMethod(fasthttp.MethodPost, c.url)
		assert.Equal(t, string(df.req.Header.Method()), fasthttp.MethodPost)
		assert.Equal(t, c.requestUrl, string(df.req.RequestURI()))
	}
}

func TestDataFlowPostDo(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		assert.Nil(t, err)
	}))
	defer s.Close()

	df := testNewDataFlow()

	var ret string
	_, err := df.newMethod(fasthttp.MethodPost, s.URL).BindString(&ret).Do(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, "ok", ret)
}

func TestDataFlowPostJson(t *testing.T) {
	src := &testJsonBody{
		A: "test",
		B: 2,
		C: 3,
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)

		var dst testJsonBody
		json.Unmarshal(b, &dst)
		assert.EqualValues(t, src, &dst)
	}))
	defer s.Close()

	df := testNewDataFlow()
	_, err := df.newMethod(fasthttp.MethodPost, s.URL).SetJson(&src).Do(context.Background())
	assert.Nil(t, err)
}

func TestDataFlowSetWWWForm(t *testing.T) {
	v := url.Values{}
	v.Set("key1", "1")
	v.Set("key2", "2")
	str := "key1=1&key2=2"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.Nil(t, err)
		assert.EqualValues(t, str, b)
	}))
	defer s.Close()

	df := testNewDataFlow()
	_, err := df.newMethod(fasthttp.MethodPost, s.URL).SetWWWForm(v).Do(context.Background())
	assert.Nil(t, err)
}

func TestDataFlowAddWWWForm(t *testing.T) {
	str := "key1=1&key1=2&key1=3&key2=4"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.Nil(t, err)
		assert.EqualValues(t, str, b)
	}))
	defer s.Close()

	df := testNewDataFlow()
	_, err := df.newMethod(fasthttp.MethodPost, s.URL).AddWWWForm("key1", "1", "2").AddWWWForm("key1", "3").AddWWWForm("key2", "4").Do(context.Background())
	assert.Nil(t, err)
}

func TestDataFlowSetBody(t *testing.T) {
	body := []byte("test body")
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.EqualValues(t, body, b)
	}))
	defer s.Close()

	df := testNewDataFlow()
	_, err := df.newMethod(fasthttp.MethodPost, s.URL).SetBody(body).Do(context.Background())
	assert.Nil(t, err)
}

func TestDataFlowTimeout(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 50)
	}))
	defer s.Close()

	df := testNewDataFlow()
	_, err := df.newMethod(fasthttp.MethodPost, s.URL).Timeout(time.Millisecond).Do(context.Background())
	assert.NotNil(t, err)

	df = testNewDataFlow()
	_, err = df.newMethod(fasthttp.MethodPost, s.URL).Do(context.Background())
	assert.Nil(t, err)
}

func TestDataFlowDegrade(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(HeaderKeyDowngrade, HeaderValueDowngradeStatic)
	}))
	defer s.Close()

	df := testNewDataFlow()
	var a int
	header := make(H)
	_, err := df.newMethod(fasthttp.MethodGet, s.URL).BindHeader(header).Degrade(func() error {
		a = 10
		return nil
	}).Do(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, HeaderValueDowngradeStatic, header.Get(HeaderKeyDowngrade))
	assert.Equal(t, 10, a)
}

func TestDataFlowCircuitBreaker(t *testing.T) {
	retCode := 500
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(retCode)
	}))
	defer s.Close()

	sentinelOptions := xsentinel.Option{
		CircuitBreakerRules: []*circuitbreaker.Rule{
			{
				Resource:         "count",
				Strategy:         circuitbreaker.ErrorCount,
				RetryTimeoutMs:   10,
				MinRequestAmount: 1,
				StatIntervalMs:   5000,
				MaxAllowedRtMs:   10,
				Threshold:        1,
			},
		},
	}

	fakeError := errors.New("fake error")
	cbFunc := func() error {
		return fakeError
	}
	err := xsentinel.Init(&sentinelOptions)
	assert.Nil(t, err)

	// 第一次出错返回
	statusCode, err := testNewDataFlow().newMethod(fasthttp.MethodGet, s.URL).CircuitBreaker("count", cbFunc).Do(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, retCode, statusCode)

	// 第二次触发熔断
	var blockError *xsentinel.BlockError
	statusCode, err = testNewDataFlow().newMethod(fasthttp.MethodGet, s.URL).CircuitBreaker("count", cbFunc).Do(context.Background())
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &(blockError)))
	assert.True(t, errors.Is(err, fakeError))
	assert.Equal(t, 0, statusCode)

	// 第三次等待后恢复
	time.Sleep(time.Millisecond * 20)
	statusCode, err = testNewDataFlow().newMethod(fasthttp.MethodGet, s.URL).CircuitBreaker("count", cbFunc).Do(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, retCode, statusCode)
}

func TestDataFlowReuse(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer s.Close()

	df := testNewDataFlow()
	tmp := df.newMethod(fasthttp.MethodGet, s.URL)
	_, err := tmp.Do(context.Background())
	assert.Nil(t, err)

	defer func() {
		e := recover()
		assert.NotNil(t, e)
	}()
	tmp.newMethod(fasthttp.MethodGet, s.URL).Do(context.Background())
}
func testNewInt(i int) *int {
	return &i
}

// type testReporter struct {
// 	spans chan *jaeger.Span
// }

// func (r *testReporter) Report(span *jaeger.Span) {
// 	select {
// 	case r.spans <- span:
// 	default:
// 	}
// }

// Close implements Close() method of Reporter by doing nothing.
// func (r *testReporter) Close() {
// }

// func TestDoTrace(t *testing.T) {
// 	tracing.Enabled()
// 	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("ok"))
// 	}))
// 	defer s.Close()
// 	tr := &testReporter{}
// 	tt, _ := jaeger.NewTracer(
// 		"test", jaeger.NewConstSampler(true), tr,
// 	)
// 	opentracing.SetGlobalTracer(tt)
// 	df := testNewDataFlow()
// 	ctx := context.Background()
// 	df.newMethod(fasthttp.MethodGet, s.URL).Do(ctx)
// 	select {
// 	case span := <-tr.spans:
// 		tags := span.Tags()
// 		assert.Equal(t, "span.type", tags["span.type"])
// 	default:
// 	}
// }
