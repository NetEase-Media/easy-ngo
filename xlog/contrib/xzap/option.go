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

package xzap

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Option 是日志配置选项
type Option struct {
	Name            string // 代表logger
	NoFile          bool   // 是否为开发模式，如果是true，只显示到标准输出，同旧的 NoFile
	Format          string
	WritableStack   bool // 是否需要打印error及以上级别的堆栈信息
	Skip            int
	WritableCaller  bool // 是否需要打印行号函数信息
	Level           string
	Path            string
	FileName        string
	PackageLevel    map[string]string // 包级别日志等级设置
	ErrlogLevel     string            // 错误日志级别，默认error
	ErrorPath       string
	MaxAge          int  // 保留旧文件的最大天数，默认7天
	MaxBackups      int  // 保留旧文件的最大个数，默认7个
	MaxSize         int  // 在进行切割之前，日志文件的最大大小（以MB为单位）默认1024
	Compress        bool // 是否压缩/归档旧文件
	packageLogLevel map[string]zapcore.Level
}

func DefaultOption() *Option {
	return &Option{
		Name:           DefaultLoggerName,
		NoFile:         true,
		Format:         formatJSON,
		WritableCaller: true,
		Skip:           2,
		WritableStack:  false,
		Level:          zap.InfoLevel.String(),
		Path:           "./logs",
		FileName:       "esay-ngo", //env.GetAppName(),
		PackageLevel:   make(map[string]string),
		ErrlogLevel:    zap.ErrorLevel.String(),
		ErrorPath:      "./logs",
		MaxAge:         7,
		MaxBackups:     7,
		MaxSize:        1024,
		Compress:       false,
	}
}

func checkOption(opt *Option) error {
	if opt.Name == "" {
		return errors.New("log name can not be nil")
	}
	if len(opt.PackageLevel) > 0 && opt.WritableCaller {
		for packageStr, packageLevelStr := range opt.PackageLevel {
			packageLevel, err := zapcore.ParseLevel(packageLevelStr)
			if err == nil {
				if opt.packageLogLevel == nil {
					opt.packageLogLevel = make(map[string]zapcore.Level)
				}
				opt.packageLogLevel[packageStr] = packageLevel
			}
		}
	}
	return nil
}
