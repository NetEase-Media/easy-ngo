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

package timeout

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func errorResponse(c *gin.Context) {
	c.String(http.StatusRequestTimeout, "timeout")
}

func doSomething(c *gin.Context) {
	q, _ := c.GetQuery("t")
	t, _ := strconv.ParseInt(q, 10, 64)
	time.Sleep(time.Millisecond * time.Duration(t))
	c.String(http.StatusOK, "success")
}

func TestTimeout(t *testing.T) {
	r := gin.New()
	r.GET("/", Timeout(WithTimeout(50*time.Millisecond), WithHandler(doSomething), WithErrorHandler(errorResponse)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?t=100", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestTimeout, w.Code)
	assert.Equal(t, "timeout", w.Body.String())
}

func TestWithoutTimeout(t *testing.T) {
	r := gin.New()
	r.GET("/", Timeout(WithTimeout(50*time.Millisecond), WithHandler(doSomething), WithErrorHandler(errorResponse)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?t=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
}
