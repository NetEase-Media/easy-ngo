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

package xfmt

import "testing"

func TestXFmt(t *testing.T) {
	xFmt, _ := Default()
	xFmt.Debugf("%s %d", "test1", 1)
	xFmt.Infof("%s %d %t", "test2", 2, true)
	xFmt.Warnf("%s %d %t %b", "test3", 3, true, 3)
	xFmt.Errorf("%s %d %t %b %x", "test4", 4, true, 4, 4)

}

func BenchmarkXFmt(b *testing.B) {
	xFmt, _ := Default()
	for i := 0; i < b.N; i++ {
		xFmt.Debugf("%s %d", "test1", 1)
		xFmt.Infof("%s %d %t", "test2", 2, true)
		xFmt.Warnf("%s %d %t %b", "test3", 3, true, 3)
		xFmt.Errorf("%s %d %t %b %x", "test4", 4, true, 4, 4)
	}
}
