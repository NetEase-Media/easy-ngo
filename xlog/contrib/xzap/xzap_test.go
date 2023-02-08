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
	"testing"
)

func TestXZap(t *testing.T) {
	xZap1, _ := New(DefaultOption())
	xZap1.Infof("%s %d", "test1", 1)
	xZap1.Warnf("%s %d %t", "test2", 2, true)
	xZap1.Errorf("%s %d %t %b", "test3", 3, true, 3)
	xZap1.Debugf("%s %d %t %b %x", "test4", 4, true, 4, 4)
}
func BenchmarkXZap(b *testing.B) {
	xZap := Default()
	for i := 0; i < b.N; i++ {
		xZap.Infof("%s %d", "test1", 1)
		xZap.Warnf("%s %d %t", "test2", 2, true)
		xZap.Errorf("%s %d %t %b", "test3", 3, true, 3)
		xZap.Debugf("%s %d %t %b %x", "test4", 4, true, 4, 4)
	}

	Option := DefaultOption()
	Option.Level = "debug"
	xZap1, _ := New(Option)
	for i := 0; i < b.N; i++ {
		xZap1.Infof("%s %d", "test1", 1)
		xZap1.Warnf("%s %d %t", "test2", 2, true)
		xZap1.Errorf("%s %d %t %b", "test3", 3, true, 3)
		xZap1.Debugf("%s %d %t %b %x", "test4", 4, true, 4, 4)
	}
}
