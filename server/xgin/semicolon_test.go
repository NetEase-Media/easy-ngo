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

package xgin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func doSomething(c *gin.Context) {
	q, _ := c.GetQuery("q")
	// log.Info(q)
	c.String(http.StatusOK, q)
}

func TestNoSemicolon(t *testing.T) {
	r := gin.New()
	r.GET("/", doSomething)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?q=a;b;c", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "a", w.Body.String())
}

func TestSemicolon(t *testing.T) {
	r := gin.New()
	r.Use(SemicolonMiddleware())
	r.GET("/", doSomething)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?q=a;b;c", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "a;b;c", w.Body.String())
}
