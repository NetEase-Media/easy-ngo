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

import "strconv"

// Code is the error code.
type Code uint32

// String returns the string representation of the status.
func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case Canceled:
		return "Canceled"
	case Unknown:
		return "Unknown"
	case InvalidArgument:
		return "InvalidArgument"
	case DeadlineExceeded:
		return "DeadlineExceeded"
	case NotFound:
		return "NotFound"
	case AlreadyExists:
		return "AlreadyExists"
	case PermissionDenied:
		return "PermissionDenied"
	case ResourceExhausted:
		return "ResourceExhausted"
	case FailedPrecondition:
		return "FailedPrecondition"
	case Aborted:
		return "Aborted"
	case OutOfRange:
		return "OutOfRange"
	case Unimplemented:
		return "Unimplemented"
	case Internal:
		return "Internal"
	case Unavailable:
		return "Unavailable"
	case DataLoss:
		return "DataLoss"
	case Unauthenticated:
		return "Unauthenticated"
	default:
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}

const (
	OK = iota
	Canceled
	Unknown
	InvalidArgument
	DeadlineExceeded
	NotFound
	AlreadyExists
	PermissionDenied
	ResourceExhausted
	FailedPrecondition
	Aborted
	OutOfRange
	Unimplemented
	Internal
	Unavailable
	DataLoss
	Unauthenticated
)

var (
	SuccessReason  = "SUCCESS"
	SuccessMessage = "success"
	UnknownReason  = "UNKNOWN"
	UnknownMessage = "unknown"
)

// SuccessStatus returns a success status.
func SuccessStatus() *Status {
	return &Status{Code: OK, Reason: SuccessReason, Message: SuccessMessage}
}

// UnknownStatus returns an unknown status.
func UnknownStatus(err error) *Status {
	msg := UnknownMessage
	if err != nil {
		msg = err.Error()
	}
	return &Status{Code: Unknown, Reason: UnknownReason, Message: msg}
}

// Status is a micro status.
type Status struct {
	Code     Code
	Reason   string
	Message  string
	Metadata map[string]string
}

// Err returns the error.
func (s *Status) Err() error {
	if s.Code == OK {
		return nil
	}
	return &Error{Status: s}
}
