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

package nlog

import (
	"log"
	"testing"

	"github.com/NetEase-Media/easy-ngo/xlog"
)

func TestNlog(t *testing.T) {

	nlog, _ := Default()
	// s, ok := nlog.(xlog.Logger)
	// if !ok {
	// 	t.Errorf("Nlog is not a xlog.Logger")
	// }
	nlog.Infof("%T", nlog)
	var _ = (xlog.Logger)(nlog)
	nlog.Debugf("%s %d", "test1", 1)
	nlog.Infof("%s %d %t", "test2", 2, true)
	nlog.Warnf("%s %d %t %b", "test3", 3, true, 3)
	nlog.Errorf("%s %d %t %b %x", "test4", 4, true, 4, 4)

}
func BenchmarkNlog(b *testing.B) {
	log1, _ := Default()
	log2 := &Nlog{flag: log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile}
	log3 := &Nlog{flag: log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC}
	log4 := &Nlog{flag: log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile}
	for i := 0; i < b.N; i++ {
		log1.Infof("%s %d", "test1", 1)
		log2.Warnf("%s %d %t", "test2", 2, true)
		log3.Errorf("%s %d %t %b", "test3", 3, true, 3)
		log4.Debugf("%s %d %t %b %x", "test4", 4, true, 4, 4)
	}
}
