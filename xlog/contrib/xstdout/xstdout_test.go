// Copyright 2023 NetEase Media Technology（Beijing）Co., Ltd.
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

package xstdout

import (
	"testing"

	"github.com/NetEase-Media/easy-ngo/xlog"
)

var logger xlog.Logger

func TestMain(m *testing.M) {
	logger = New()
	m.Run()
}

func TestDebugf(t *testing.T) {
	logger.Debugf("test %s", "debug")
}

func TestInfof(t *testing.T) {
	logger.Infof("test %s", "info")
}

func TestWarnf(t *testing.T) {
	logger.Warnf("test %s", "warn")
}

func TestErrorf(t *testing.T) {
	logger.Errorf("test %s", "error")
}

func TestFatalf(t *testing.T) {
	logger.Fatalf("test %s", "fatal")
}

func TestPanicf(t *testing.T) {
	logger.Panicf("test %s", "panic")
}
