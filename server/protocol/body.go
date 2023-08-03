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

package protocol

import (
	"errors"
	"net/http"
)

// HttpBody 是写入http body的json数据结构
type HttpBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// GetError 判断回复是否错误，如果是则返回对应错误对象
func (b *HttpBody) GetError() *Error {
	if b.Code == 0 {
		return nil
	}

	return &Error{
		Code: b.Code,
		Err:  errors.New(b.Message),
	}
}

// JsonBody 生成成功回复的http code和body
func JsonBody(data interface{}) (statusCode int, body *HttpBody) {
	return Success(data)
}

// Success 业务处理成功
func Success(data interface{}) (statusCode int, body *HttpBody) {
	return Response(http.StatusOK, 0, "成功", data)
}

// Result 业务处理结果
func Result(code int, message string, data interface{}) (statusCode int, body *HttpBody) {
	return Response(http.StatusOK, code, message, data)
}

// Response 通用业务处理
func Response(status int, code int, message string, data interface{}) (statusCode int, body *HttpBody) {
	return status, &HttpBody{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
