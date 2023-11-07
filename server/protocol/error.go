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
	"fmt"
	"net/http"
)

// ErrorJsonBody 生成错误信息的http code和body
func ErrorJsonBody(errorCode int) (int, *HttpBody) {
	statusCode, ok := errorStatus[errorCode]
	if !ok {
		statusCode = http.StatusOK // 理论上说未知错误应该返回500，但是程序这样写了，没办法🤷‍
	}
	return statusCode, &HttpBody{
		Code:    errorCode,
		Message: errorMessages[errorCode],
	}
}

// Fail 业务处理失败
func Fail(code int, message string) (statusCode int, body *HttpBody) {
	return http.StatusOK, &HttpBody{
		Code:    code,
		Message: message,
	}
}

const (
	SystemError        = 1000000
	DBError            = 1000001
	CacheError         = 1000002
	ThirdServiceError  = 1000003
	ParamsLost         = 1000100
	ParamsNotValid     = 1000101
	ResouceNotExist    = 1000102
	DataOutOfThreshold = 1000103
	FrequentOpration   = 1000104
	RepeatOpration     = 1000105
	IllegalRequest     = 1000106
	DataHasExists      = 1000107
	PermissionDenied   = 1000108
	AntiCheating       = 1000109
	UnsupportClient    = 1000110
	UnsupportOS        = 1000111
	AccountFrozen      = 1000200
	AccountLock        = 1000201
	TokenError         = 1000202
)

var errorMessages = map[int]string{
	SystemError:        "服务器内部错误",
	DBError:            "服务器内部错误",
	CacheError:         "服务器内部错误",
	ThirdServiceError:  "服务器内部错误",
	ParamsLost:         "请求参数缺失",
	ParamsNotValid:     "存在不合法的请求参数",
	ResouceNotExist:    "资源不存在",
	DataOutOfThreshold: "数据超过阈值",
	FrequentOpration:   "操作频繁",
	RepeatOpration:     "重复操作",
	IllegalRequest:     "非法请求",
	DataHasExists:      "数据已存在",
	PermissionDenied:   "权限不足",
	AntiCheating:       "请求被拦截",
	UnsupportClient:    "不支持的客户端",
	UnsupportOS:        "不支持的操作系统",
	AccountFrozen:      "账号异常-需打开安全中心申诉",
	AccountLock:        "账号异常-需打开安全中心解锁",
	TokenError:         "token校验失败",
}

var errorStatus = map[int]int{
	SystemError:       http.StatusInternalServerError,
	DBError:           http.StatusInternalServerError,
	CacheError:        http.StatusInternalServerError,
	ThirdServiceError: http.StatusInternalServerError,
	ParamsLost:        http.StatusBadRequest,
	ParamsNotValid:    http.StatusBadRequest,
	TokenError:        http.StatusBadRequest,
}

// Error 用来将运行错误包装成标准协议的错误
type Error struct {
	Code int
	Err  error
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	return fmt.Sprintf("code:%d, error:%s, message:%s", e.Code, errorMessages[e.Code], e.Err.Error())
}

func (e *Error) HttpBody() (int, *HttpBody) {
	statusCode, body := ErrorJsonBody(e.Code)
	body.Data = e.Err.Error()
	return statusCode, body
}
