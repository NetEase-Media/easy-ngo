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
	"fmt"

	"github.com/alibaba/sentinel-golang/core/base"
)

// BlockError 用来存储sentinel的熔断错误和用户自身错误
type BlockError struct {
	BlockErr *base.BlockError
	Err      error
}

func (e *BlockError) Unwrap() error {
	return e.Err
}

func (e *BlockError) Error() string {
	return fmt.Sprintf("sentinel block error: %s, wrapped error: %v", e.BlockErr.Error(), e.Err)
}
