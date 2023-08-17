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
	"context"
	"encoding/json"
	"fmt"
	"net/textproto"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/djimenez/iconv-go"
	"github.com/valyala/fasthttp"
)

type bodyType uint8

const (
	typeNil bodyType = iota
	typeInt
	typeFloat
	typeString
	typeBytes
	typeJson
)

var (
	charsetReg = regexp.MustCompile(`charset=([^;]+)`)
)

// http header类型
type H = textproto.MIMEHeader

// url query结构
type Query = url.Values

// x-www-form-urlencoded body类型
type WWWForm = url.Values

type DoFunc func(*DataFlow, context.Context) (int, error)

// DataFlow 是核心数据结构，用来保存http请求的中间状态数据，并负责发送请求和解析回复。
type DataFlow struct {
	// req 保存请求的中间属性。在Do执行后内存会被释放，不可再使用！
	req    *fasthttp.Request
	client *fasthttp.Client
	header H
	query  Query

	wwwForm         WWWForm       // 使用AddWWWFrom或SetWWWForm写入的数据
	timeout         time.Duration // 单次请求的超时时间
	degradeCallback func() error  // 降级回调函数
	cbCallback      func() error  // 熔断回调

	// 绑定回复的http body
	// TODO: 因为只可能使用其中一种，可以考虑用interface保存，使用时再转换
	bodyBindType bodyType
	bodyInt      *int
	bodyFloat    *float64
	bodyString   *string
	bodyBytes    *[]byte
	bodyJson     interface{}

	// 绑定请求的http header
	headerBinder H

	Err error

	do1 DoFunc
}

func newDataFlow(c *Xfasthttp) *DataFlow {
	df := &DataFlow{
		client: c.client,
	}
	df.do1 = func(df1 *DataFlow, c context.Context) (int, error) {
		return df1.doInternal()
	}
	return df
}

func (df *DataFlow) WrapDoFunc(f func(DoFunc) DoFunc) {
	df.do1 = f(df.do1)
}

// newMethod 以method和url初始化请求
func (df *DataFlow) newMethod(method, url string) *DataFlow {
	df.req = fasthttp.AcquireRequest()
	df.req.Header.SetMethod(method)
	df.req.SetRequestURI(url)
	return df
}

// Timeout 设置请求的超时时间
func (df *DataFlow) Timeout(t time.Duration) *DataFlow {
	df.timeout = t
	return df
}

// SetQuery 设置url query字段
func (df *DataFlow) SetQuery(q Query) *DataFlow {
	df.query = q
	return df
}

// AddQuery 增加一个url query键值对
func (df *DataFlow) AddQuery(key, value string) *DataFlow {
	if df.query == nil {
		df.query = make(Query)
	}
	df.query.Add(key, value)
	return df
}

// SetContentType 设置http header的ContentType
func (df *DataFlow) SetContentType(t string) *DataFlow {
	df.req.Header.SetContentType(t)
	return df
}

// SetBody 直接设置http body
func (df *DataFlow) SetBody(b []byte) *DataFlow {
	df.req.SetBody(b)
	return df
}

// SetJson 设置请求体，格式是json
func (df *DataFlow) SetJson(body interface{}) *DataFlow {
	b, err := json.Marshal(body)
	if err != nil {
		df.Err = err
		xlog.Errorf("encoding body failed: %s", err.Error())
		return df
	}

	// log.Infof("%s", string(b))
	df.req.SetBody(b)
	df.req.Header.SetContentType("application/json")
	return df
}

// SetWWWForm 设置请求体，格式是x-www-form-urlencoded
func (df *DataFlow) SetWWWForm(f WWWForm) *DataFlow {
	df.wwwForm = f
	return df
}

func (df *DataFlow) SetFormWithMap(data map[string]interface{}) *DataFlow {
	if df.wwwForm == nil {
		df.wwwForm = make(WWWForm)
	}
	for k, v := range data {
		df.wwwForm.Add(k, fmt.Sprintf("%v", v))
	}
	return df
}

// AddWWWForm 提供简便的x-www-form-urlencoded格式调用接口
func (df *DataFlow) AddWWWForm(key string, values ...string) *DataFlow {
	if df.wwwForm == nil {
		df.wwwForm = make(WWWForm)
	}
	for _, value := range values {
		df.wwwForm.Add(key, value)
	}
	return df
}

// AddHeader 设置请求的http header，传入map结构
func (df *DataFlow) AddHeader(h H) *DataFlow {
	if h == nil {
		return df
	}

	for k, arr := range h {
		for _, v := range arr {
			df.req.Header.Add(k, v)
		}
	}
	return df
}

// AddHeaderKV 设置请求的http header
func (df *DataFlow) AddHeaderKV(key string, values ...string) *DataFlow {
	for _, value := range values {
		df.req.Header.Add(key, value)
	}
	return df
}

// BindJson 将json结构对象与body绑定
// 注意obj必须是一个指针
func (df *DataFlow) BindJson(obj interface{}) *DataFlow {
	df.bodyBindType = typeJson
	df.bodyJson = obj
	return df
}

// BindString 将string值与body绑定
func (df *DataFlow) BindString(s *string) *DataFlow {
	df.bodyBindType = typeString
	df.bodyString = s
	return df
}

// BindFloat 将float值与body绑定
func (df *DataFlow) BindFloat(f *float64) *DataFlow {
	df.bodyBindType = typeFloat
	df.bodyFloat = f
	return df
}

// BindBytes 将bytes值与body绑定
func (df *DataFlow) BindBytes(b *[]byte) *DataFlow {
	df.bodyBindType = typeBytes
	df.bodyBytes = b
	return df
}

// BindInt 将int值与body绑定
func (df *DataFlow) BindInt(i *int) *DataFlow {
	df.bodyBindType = typeInt
	df.bodyInt = i
	return df
}

// processRequest 将请求缓存解析并写入到request中
func (df *DataFlow) processRequest() {
	// 注意www-form的优先级大于其他
	if df.wwwForm != nil {
		df.req.SetBodyString(df.wwwForm.Encode())
		df.req.Header.SetContentType("application/x-www-form-urlencoded")
	}

	// 将url query写入request
	if df.query != nil {
		df.req.URI().SetQueryString(df.query.Encode())
	}
}

// processResponse 解析回复，存储到绑定的变量中
func (df *DataFlow) processResponse(res *fasthttp.Response) error {
	xlog.Debugf("http recv response header\n%s\n body\n%s", &res.Header, string(res.Body()))

	if err := df.encodeHeader(&res.Header); err != nil {
		return err
	}

	// 处理降级
	if df.header.Get(HeaderKeyDowngrade) != "" && df.degradeCallback != nil {
		return df.degradeCallback()
	}

	// 如果是5XX则记录熔断错误

	if err := df.encodeBody(res); err != nil {
		return err
	}

	return nil
}

// send 根据当前状态选择发送请求
func (df *DataFlow) send(res *fasthttp.Response) (err error) {
	xlog.Debugf("http send request header {%s} body {%s}", df.req.Header.String(), string(df.req.Body()))

	if df.timeout != 0 {
		return df.client.DoTimeout(df.req, res, df.timeout)
	}

	// 回复可能有重定向，所以使用DoRedirects发送请求
	return df.client.DoRedirects(df.req, res, defaultMaxRedirectsCount)
}

// Do 调用fasthttp client发送请求，并解析回复数据
// 注意一旦使用后DataFlow将不可再使用
func (df *DataFlow) doInternal() (statusCode int, err error) {
	// var stats *httpclient.StatsHolder
	defer func() {
		df.release()
		// 清空，防止垃圾和重复使用
		df.reset()
	}()
	if df.Err != nil {
		return 0, df.Err
	}

	df.processRequest()
	res := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseResponse(res)
	}()

	if err = df.send(res); err != nil {
		return
	}

	statusCode = res.StatusCode()
	err = df.processResponse(res)
	return
}

func (df *DataFlow) Do(ctx context.Context) (statusCode int, err error) {
	return df.do1(df, ctx)
}

// reset 清理对象，防止重复使用。如果使用sync.Pool必须调用。
func (df *DataFlow) reset() {
	df.req = nil
	df.client = nil
	df.header = nil
	df.query = nil
	df.wwwForm = nil
	df.headerBinder = nil
	df.degradeCallback = nil
	df.cbCallback = nil
	df.Err = nil
}

// Degrade 注册降级回调函数
func (df *DataFlow) Degrade(f func() error) *DataFlow {
	df.degradeCallback = f
	return df
}

// encodeBody 将回复中的body写入绑定的对象中
func (df *DataFlow) encodeBody(res *fasthttp.Response) (err error) {
	body := res.Body()
	switch df.bodyBindType {
	case typeNil:
		return

	case typeInt:
		if df.bodyInt == nil {
			err = fmt.Errorf("empty int binder")
			break
		}
		var i int
		i, err = strconv.Atoi(string(body))
		if err != nil {
			break
		}
		*df.bodyInt = i

	case typeFloat:
		if df.bodyFloat == nil {
			err = fmt.Errorf("empty float binder")
			break
		}
		var f float64
		f, err = strconv.ParseFloat(string(body), 64)
		if err != nil {
			break
		}
		*df.bodyFloat = f

	case typeString:
		if df.bodyString == nil {
			err = fmt.Errorf("empty string binder")
			break
		}

		if charset := getCharset(string(res.Header.ContentType())); !strings.EqualFold(charset, "utf-8") {
			output, e := iconv.ConvertString(string(body), charset, "utf-8")
			if e != nil {
				xlog.Errorf("convert string %s from %s to %s error: %v", charset, "utf-8", e)
				*df.bodyString = string(body)
				return
			}
			*df.bodyString = output
		} else {
			*df.bodyString = string(body)
		}

	case typeBytes:
		if df.bodyBytes == nil {
			err = fmt.Errorf("empty bytes binder")
			break
		}

		// 因为response在Do之后会释放，一定要把body复制出来
		if len(*df.bodyBytes) < len(body) {
			*df.bodyBytes = make([]byte, len(body))
		}
		copy(*df.bodyBytes, body)

	case typeJson:
		if res.Header.ContentType() == nil {
			*df.bodyString = string(body)
		}

		if charset := getCharset(string(res.Header.ContentType())); !strings.EqualFold(charset, "utf-8") {
			output, e := iconv.ConvertString(string(body), charset, "utf-8")
			if e != nil {
				xlog.Errorf("convert string %s from %s to %s error: %v", charset, "utf-8", e)
				err = json.Unmarshal(body, df.bodyJson)
				return
			}
			err = json.Unmarshal([]byte(output), df.bodyJson)
		} else {
			err = json.Unmarshal(body, df.bodyJson)
		}
	default:
		panic(fmt.Sprintf("wrong bind type %d", df.bodyBindType))
	}

	return
}

// getCharset 从contentType中获取编码
func getCharset(contentType string) string {
	contentType = strings.TrimSpace(strings.ToLower(contentType))
	// 捕获编码
	re := charsetReg.FindAllStringSubmatch(contentType, 1)
	if len(re) > 0 {
		c := re[0][1]
		return c
	}
	return "utf-8"
}

func (df *DataFlow) BindHeader(header H) *DataFlow {
	df.headerBinder = header
	return df
}

// encodeHeader 将回复中的header写入绑定的对象中
func (df *DataFlow) encodeHeader(header *fasthttp.ResponseHeader) (err error) {
	// 绑定的header与请求的缓存指向同一map
	if df.headerBinder != nil {
		df.header = df.headerBinder
	} else {
		df.header = make(H)
	}

	header.VisitAll(func(key, value []byte) {
		sk := string(key)
		df.header[sk] = append(df.header[sk], string(value))
	})

	return
}
