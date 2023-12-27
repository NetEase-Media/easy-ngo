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

// if not config zap logger, use stdout logger
package xstdout

import "fmt"

type XStdout struct {
}

func New() *XStdout {
	return &XStdout{}
}

func (x *XStdout) Debugf(msg string, params ...interface{}) {
	withPrefix := fmt.Sprintf("[DEBUG] %s \n", msg)
	fmt.Printf(withPrefix, params...)
}

func (x *XStdout) Infof(msg string, params ...interface{}) {
	withPrefix := fmt.Sprintf("[INFO] %s \n", msg)
	fmt.Printf(withPrefix, params...)
}

func (x *XStdout) Errorf(msg string, params ...interface{}) {
	withPrefix := fmt.Sprintf("[ERROR] %s \n", msg)
	fmt.Printf(withPrefix, params...)
}

func (x *XStdout) Warnf(msg string, params ...interface{}) {
	withPrefix := fmt.Sprintf("[WARN] %s \n", msg)
	fmt.Printf(withPrefix, params...)
}

func (x *XStdout) Fatalf(msg string, params ...interface{}) {
	withPrefix := fmt.Sprintf("[FATAL] %s \n", msg)
	fmt.Printf(withPrefix, params...)
}

func (x *XStdout) Panicf(msg string, params ...interface{}) {
	withPrefix := fmt.Sprintf("[PANIC] %s \n", msg)
	fmt.Printf(withPrefix, params...)
}
