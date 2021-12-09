package gcode

import (
	"github.com/RussellLuo/kun/pkg/werror"
)

// The following error codes are borrowed from gRPC, see
// https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto

var (
	ErrInvalidArgument    = werror.Wrapf(nil, "InvalidArgument")    // HTTP Mapping: 400
	ErrFailedPrecondition = werror.Wrapf(nil, "FailedPrecondition") // HTTP Mapping: 400
	ErrOutOfRange         = werror.Wrapf(nil, "OutOfRange")         // HTTP Mapping: 400
	ErrUnauthenticated    = werror.Wrapf(nil, "Unauthenticated")    // HTTP Mapping: 401
	ErrPermissionDenied   = werror.Wrapf(nil, "PermissionDenied")   // HTTP Mapping: 403
	ErrNotFound           = werror.Wrapf(nil, "NotFound")           // HTTP Mapping: 404
	ErrAborted            = werror.Wrapf(nil, "Aborted")            // HTTP Mapping: 409
	ErrAlreadyExists      = werror.Wrapf(nil, "AlreadyExists")      // HTTP Mapping: 409
	ErrResourceExhausted  = werror.Wrapf(nil, "ResourceExhausted")  // HTTP Mapping: 429
	ErrCancelled          = werror.Wrapf(nil, "Cancelled")          // HTTP Mapping: 499
	ErrDataLoss           = werror.Wrapf(nil, "DataLoss")           // HTTP Mapping: 500
	ErrUnknown            = werror.Wrapf(nil, "Unknown")            // HTTP Mapping: 500
	ErrInternal           = werror.Wrapf(nil, "Internal")           // HTTP Mapping: 500
	ErrNotImplemented     = werror.Wrapf(nil, "NotImplemented")     // HTTP Mapping: 501
	ErrUnavailable        = werror.Wrapf(nil, "Unavailable")        // HTTP Mapping: 503
	ErrDeadlineExceeded   = werror.Wrapf(nil, "DeadlineExceeded")   // HTTP Mapping: 504
)
