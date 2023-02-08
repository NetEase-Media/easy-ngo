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

package grpc

import (
	"fmt"
	"testing"

	"github.com/NetEase-Media/easy-ngo/microservices/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFromRPCError(t *testing.T) {
	var err, wantErr error
	err = FromRPCError(nil)
	assert.NoError(t, err)

	err = FromRPCError(fmt.Errorf("unknown"))
	wantErr = errors.UnknownError(fmt.Errorf("unknown"))
	assert.True(t, err.(*errors.Error).Is(wantErr))

	err = FromRPCError(status.Errorf(codes.Unknown, "test"))
	wantErr = errors.New(errors.Unknown, errors.UnknownReason, "test")
	assert.True(t, err.(*errors.Error).Is(wantErr))
}

func TestToRPCError(t *testing.T) {
	var err, wantErr error
	err = ToRPCError(nil)
	assert.NoError(t, err)

	err = ToRPCError(fmt.Errorf("unknown"))
	wantErr = status.Errorf(codes.Unknown, fmt.Errorf("unknown").Error())
	assert.True(t, status.Convert(err).String() == status.Convert(wantErr).String())

	err = ToRPCError(errors.New(errors.Unknown, errors.UnknownReason, "test"))
	wantErr = status.Errorf(codes.Unknown, "test")
	assert.True(t, status.Convert(err).String() == status.Convert(wantErr).String())
}
