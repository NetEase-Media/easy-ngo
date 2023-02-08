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

package errors

import (
	"errors"
	"fmt"
)

// UnknownError returns an error with unknown status.
func UnknownError(err error) *Error {
	return UnknownStatus(err).Err().(*Error)
}

// New returns an error with the specified code, reason, message and metadata.
func New(c Code, reason, msg string, kvs ...string) *Error {
	var md map[string]string
	if len(kvs) > 0 {
		if len(kvs)%2 != 0 {
			kvs = append(kvs, "")
		}
		md = make(map[string]string, len(kvs)/2)
		for i := 0; i < len(kvs); i += 2 {
			md[kvs[i]] = kvs[i+1]
		}
	}
	return &Error{&Status{Code: c, Reason: reason, Message: msg, Metadata: md}}
}

// Newf returns an error with the specified code, reason and message.
func Newf(c Code, reason, format string, a ...interface{}) *Error {
	return New(c, reason, fmt.Sprintf(format, a...))
}

// Error is a micro error.
type Error struct {
	*Status
}

// Error returns the error string.
func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, reason: %s, msg: %s, md: %v", e.Code, e.Reason, e.Message, e.Metadata)
}

// WithMetadata returns an error with metadata.
func (e *Error) WithMetadata(md map[string]string) *Error {
	status := Status{
		Code:     e.Code,
		Reason:   e.Reason,
		Message:  e.Message,
		Metadata: md,
	}
	return &Error{&status}
}

// Is returns true if the error is of the specified type.
func (e *Error) Is(err error) bool {
	if er := new(Error); errors.As(err, &er) {
		return er.Code == e.Code && er.Reason == e.Reason
	}
	return false
}

// FromError returns a status from an error.
func FromError(err error) *Status {
	if err == nil {
		return SuccessStatus()
	}

	if se := new(Error); errors.As(err, &se) {
		return se.Status
	}

	return UnknownStatus(err)
}
