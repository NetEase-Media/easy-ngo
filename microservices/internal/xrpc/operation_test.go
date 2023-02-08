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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseToOperation(t *testing.T) {
	var op Operation
	var ok bool

	op, ok = ParseToOperation("/xx")
	assert.False(t, ok)

	op, ok = ParseToOperation("/pkg.service/method")
	assert.True(t, ok)
	assert.Equal(t, "pkg", op.Pkg)
	assert.Equal(t, "service", op.Service)
	assert.Equal(t, "method", op.Method)
}
