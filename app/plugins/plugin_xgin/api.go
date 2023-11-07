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

package pluginxgin

import (
	"github.com/NetEase-Media/easy-ngo/server/contrib/xgin"
	"github.com/gin-gonic/gin"
)

var s *xgin.Server

func WithServer(s1 *xgin.Server) {
	s = s1
}

func GetServer() *xgin.Server {
	return s
}

func PUT(relativePath string, handler gin.HandlerFunc) error {
	s.Engine.PUT(relativePath, handler)
	return nil
}

func GET(relativePath string, handler gin.HandlerFunc) error {
	s.Engine.GET(relativePath, handler)
	return nil
}

func POST(relativePath string, handler gin.HandlerFunc) error {
	s.Engine.POST(relativePath, handler)
	return nil
}

func DELETE(relativePath string, handler gin.HandlerFunc) error {
	s.Engine.DELETE(relativePath, handler)
	return nil
}

func PATCH(relativePath string, handler gin.HandlerFunc) error {
	s.Engine.PATCH(relativePath, handler)
	return nil
}

func HEAD(relativePath string, handler gin.HandlerFunc) error {
	s.Engine.HEAD(relativePath, handler)
	return nil
}

func OPTIONS(relativePath string, handler gin.HandlerFunc) error {
	s.Engine.OPTIONS(relativePath, handler)
	return nil
}
