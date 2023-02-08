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
	goerr "errors"

	"github.com/NetEase-Media/easy-ngo/microservices/errors"
	"github.com/NetEase-Media/easy-ngoservices/transport"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ transport.ErrorConverter = (*ErrorConverter)(nil)
var defaultErrorConverter = &ErrorConverter{}

func FromRPCError(err error) error {
	return defaultErrorConverter.FromRPCError(err)
}

func ToRPCError(err error) error {
	return defaultErrorConverter.ToRPCError(err)
}

type ErrorConverter struct {
}

func (e *ErrorConverter) FromRPCError(err error) error {
	if err == nil {
		return nil
	}
	s, _ := status.FromError(err)
	reason := errors.UnknownReason
	var metadata map[string]string
	for _, detail := range s.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			reason = d.Reason
			metadata = d.Metadata
		}
	}
	return errors.New(e.fromRPCCode(s.Code()), reason, s.Message()).WithMetadata(metadata)
}

func (e *ErrorConverter) ToRPCError(err error) error {
	if err == nil {
		return nil
	}
	if er := new(errors.Error); goerr.As(err, &er) {
		s, _ := status.New(e.toRPCCode(er.Code), er.Message).
			WithDetails(&errdetails.ErrorInfo{
				Reason:   er.Reason,
				Metadata: er.Metadata,
			})
		return s.Err()
	}
	return status.Error(codes.Unknown, err.Error())
}

func (e *ErrorConverter) toRPCCode(c errors.Code) codes.Code {
	switch c {
	case errors.OK:
		return codes.OK
	case errors.Canceled:
		return codes.Canceled
	case errors.Unknown:
		return codes.Unknown
	case errors.InvalidArgument:
		return codes.InvalidArgument
	case errors.DeadlineExceeded:
		return codes.DeadlineExceeded
	case errors.NotFound:
		return codes.NotFound
	case errors.AlreadyExists:
		return codes.AlreadyExists
	case errors.PermissionDenied:
		return codes.PermissionDenied
	case errors.ResourceExhausted:
		return codes.ResourceExhausted
	case errors.FailedPrecondition:
		return codes.FailedPrecondition
	case errors.Aborted:
		return codes.Aborted
	case errors.OutOfRange:
		return codes.OutOfRange
	case errors.Unimplemented:
		return codes.Unimplemented
	case errors.Internal:
		return codes.Internal
	case errors.Unavailable:
		return codes.Unavailable
	case errors.DataLoss:
		return codes.DataLoss
	case errors.Unauthenticated:
		return codes.Unauthenticated
	default:
		return codes.Unknown
	}
}

func (e *ErrorConverter) fromRPCCode(c codes.Code) errors.Code {
	switch c {
	case codes.OK:
		return errors.OK
	case errors.Canceled:
		return errors.Canceled
	case errors.Unknown:
		return errors.Unknown
	case errors.InvalidArgument:
		return errors.InvalidArgument
	case errors.DeadlineExceeded:
		return errors.DeadlineExceeded
	case errors.NotFound:
		return errors.NotFound
	case errors.AlreadyExists:
		return errors.AlreadyExists
	case errors.PermissionDenied:
		return errors.PermissionDenied
	case errors.ResourceExhausted:
		return errors.ResourceExhausted
	case errors.FailedPrecondition:
		return errors.FailedPrecondition
	case errors.Aborted:
		return errors.Aborted
	case errors.OutOfRange:
		return errors.OutOfRange
	case errors.Unimplemented:
		return errors.Unimplemented
	case errors.Internal:
		return errors.Internal
	case errors.Unavailable:
		return errors.Unavailable
	case errors.DataLoss:
		return errors.DataLoss
	case errors.Unauthenticated:
		return errors.Unauthenticated
	default:
		return errors.Unknown
	}
}
