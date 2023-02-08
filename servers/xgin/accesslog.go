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
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/NetEase-Media/easy-ngo/servers/xgin/accesslog"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

type AccessLogMwOption struct {
	Enabled    bool
	Pattern    string
	Path       string
	FileName   string
	NoFile     bool
	MaxAge     int  // 保留旧文件的最大天数，默认7天
	MaxBackups int  // 保留旧文件的最大个数，默认7个
	MaxSize    int  // 在进行切割之前，日志文件的最大大小（以MB为单位）默认1024
	Compress   bool // 是否压缩/归档旧文件
}

func NewDefaultAccessLogOptions() *AccessLogMwOption {
	return &AccessLogMwOption{
		Enabled:    true,
		Pattern:    accesslog.ApacheCombinedLogFormat,
		Path:       "",
		FileName:   "access.log",
		NoFile:     true,
		MaxAge:     7,
		MaxBackups: 7,
		MaxSize:    1024,
		Compress:   false,
	}
}

func AccessLogMiddleware(opt *AccessLogMwOption) gin.HandlerFunc {
	if opt == nil {
		opt = NewDefaultAccessLogOptions()
	}
	if opt.Enabled {
		if opt.NoFile {
			return accesslog.FormatWith(opt.Pattern, accesslog.WithOutput(os.Stdout))
		}

		writer, err := newRotateLog(opt)
		if err != nil {
			panic(err)
		}
		return accesslog.FormatWith(opt.Pattern, accesslog.WithOutput(writer))
	}
	return func(c *gin.Context) {
		c.Next()
	}
}

func newRotateLog(opt *AccessLogMwOption) (io.Writer, error) {
	dir, err := filepath.Abs(opt.Path)
	if err != nil {
		return nil, err
	}
	linkName := path.Join(dir, opt.FileName)

	return &lumberjack.Logger{
		Filename:   linkName,
		MaxSize:    opt.MaxSize, // megabytes
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge, //days
		LocalTime:  true,
		Compress:   opt.Compress, // disabled by default
	}, nil
}
