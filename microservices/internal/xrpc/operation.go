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

package xrpc

import (
	"fmt"
	"strings"
)

// Operation is a struct that represents a service method.
type Operation struct {
	Pkg     string
	Service string
	Method  string
}

// FullService returns the full service name.
func (o Operation) FullService() string {
	return fmt.Sprintf("%s.%s", o.Pkg, o.Service)
}

// String returns the string representation of the operation.
func (o Operation) String() string {
	return fmt.Sprintf("/%s.%s/%s", o.Pkg, o.Service, o.Method)
}

// ParseToOperation parses the string to operation.
func ParseToOperation(fullMethod string) (op Operation, ok bool) {
	str := strings.TrimLeft(fullMethod, "/")
	arr := strings.SplitN(str, "/", 2)
	if len(arr) == 2 { //nolint:gomnd
		arr2 := strings.SplitN(arr[0], ".", 2)
		if len(arr2) == 2 { //nolint:gomnd
			op.Pkg, op.Service = arr2[0], arr2[1]
			op.Method = arr[1]
			return op, true
		}
	}
	return op, false
}
