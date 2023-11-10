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
	"fmt"
	"os"
	"testing"

	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/stretchr/testify/assert"
)

const (
	_path       = "./logs"
	_error_path = "./logs"
)

func TestMain(m *testing.M) {

	// init
	setup()

	m.Run()

	// shutdown
	teardown()
}

func clearLogs() error {
	_, err := os.Stat(_path)
	if err == nil {
		err = os.RemoveAll(_path)
		return err
	}
	_, err = os.Stat(_error_path)
	if err == nil {
		err = os.RemoveAll(_error_path)
		return err
	}
	return nil
}

func setup() error {
	// remove old logs
	return clearLogs()
}

func teardown() error {
	// clean
	return clearLogs()
}

func TestXzapInfo(t *testing.T) {
	c := DefaultConfig()
	c.Path = _path
	xzap, _ := New(c)
	xlog.WithVendor(xzap)
	xlog.Infof("info %s", "test")
	fpath := fmt.Sprintf("%s/%s.%s", c.Path, c.FileName, c.Suffix)
	_, err := os.Stat(fpath)
	assert.Nil(t, err, "can not find log file")
}

func TestXzapError(t *testing.T) {
	c := DefaultConfig()
	c.ErrorPath = _error_path
	xzap, _ := New(c)
	xlog.WithVendor(xzap)
	xlog.Errorf("error %s", "test")
	fpath := fmt.Sprintf("%s/%s.%s", c.ErrorPath, c.FileName, c.ErrorSuffix)
	_, err := os.Stat(fpath)
	assert.Nil(t, err, "can not find error log file")
}
