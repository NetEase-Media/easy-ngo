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

package file

import "testing"

func TestDecode(t *testing.T) {
	// tomlBlob := `
	// 	# Some comments.
	// 	[alpha]
	// 		ip = "10.0.0.1"
	// 	[alpha.config]
	// 		Ports = [ 8001, 8002 ]
	// 		Location = "Toronto"
	// 		Created = 1987-07-05T05:45:00Z
	// 	[beta]
	// 		ip = "10.0.0.2"
	// 		[beta.config]
	// 		Ports = [ 9001, 9002 ]
	// 		Location = "New Jersey"
	// 		Created = 1887-01-05T05:55:00Z
	// `

	// ret := decodeToml(tomlBlob)
	// t.Log(ret["alpha.ip"])
}
