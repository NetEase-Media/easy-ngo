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
	"testing"
	"time"

	"github.com/NetEase-Media/easy-ngo/servers/xgin/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAccessLog(t *testing.T) {
	opt := NewDefaultAccessLogOptions()
	opt.Pattern = `%a %A %b %B %h %H %l %m %p %q %r %s %S %t %u %U %v %D %T %I %{X-Real-Ip}i %{User-Agent}i %{Content-Type}o %{xxx}c %{data}r`
	opt.NoFile = false
	r := gin.Default()
	r.Use(AccessLogMiddleware(opt))
	r.GET("/ping", func(c *gin.Context) {
		c.Set("data", "data...")
		time.Sleep(3 * time.Millisecond)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	w := util.PerformRequest(r, "GET", "/ping?p=aaaa&q=bbbb", util.Header{Key: "X-Real-IP", Value: "1.1.1.1"},
		util.Header{Key: "User-Agent", Value: "AHC/2.1 ..."})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAccessLog_branchCov(t *testing.T) {
	r := gin.Default()
	r.Use(AccessLogMiddleware(nil))
	opt := AccessLogMwOption{
		Enabled: false,
		Pattern: `%a %A %b %B %h %H %l %m %p %q %r %s %S %t %u %U %v %D %T %I %{X-Real-Ip}i %{User-Agent}i %{Content-Type}o %{xxx}c %{data}r`,
	}
	r.Use(AccessLogMiddleware(&opt))
	opt = AccessLogMwOption{
		Enabled: true,
		Pattern: `%a %A %b %B %h %H %l %m %p %q %r %s %S %t %u %U %v %D %T %I %{X-Real-Ip}i %{User-Agent}i %{Content-Type}o %{xxx}c %{data}r`,
		NoFile:  false,
	}
	r.Use(AccessLogMiddleware(&opt))
}

func TestAccessLog_(t *testing.T) {
	opt := NewDefaultAccessLogOptions()
	opt.Pattern = `%a %A %b %B %h %H %l %m %p %q %r %s %S %t %u %U %v %D %T %I %{X-Real-Ip}i %{User-Agent}i %{Content-Type}o %{xxx}c %{data}r`
	opt.NoFile = false
	r := gin.Default()
	r.Use(AccessLogMiddleware(opt))
}

func BenchmarkAccessLog(b *testing.B) {
	opt := NewDefaultAccessLogOptions()
	opt.Pattern = "common"
	opt.NoFile = false
	r := gin.Default()
	r.Use(AccessLogMiddleware(opt))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	util.RunRequest(b, r, "GET", "/ping?p=aaaa", util.Header{Key: "X-Real-IP", Value: "1.1.1.1"},
		util.Header{Key: "User-Agent", Value: "AHC/2.1"})
}
