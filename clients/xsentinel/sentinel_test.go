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

package xsentinel

import (
	"testing"

	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/stretchr/testify/assert"
)

func TestSentinel(t *testing.T) {
	opt := Option{
		FlowRules: []*flow.Rule{
			{
				Resource:               "some-test",
				TokenCalculateStrategy: flow.Direct,
				ControlBehavior:        flow.Reject,
				Threshold:              1,
				StatIntervalInMs:       10000,
			},
		},
	}
	Init(&opt)
	var succ, fail int
	for i := 0; i < 2; i++ {
		e, b := Entry("some-test")
		if b != nil {
			fail++
		} else {
			succ++
			e.Exit()
		}
	}
	assert.Equal(t, succ, 1)
	assert.Equal(t, fail, 1)
}
