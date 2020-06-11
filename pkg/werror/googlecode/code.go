package googlecode

import (
	"github.com/RussellLuo/kok/pkg/werror"
)

// The following error codes are borrowed from
// https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto

var (
	ErrInvalidArgument    = werror.Wrap(nil).SetErrorf("InvalidArgument")    // HTTP Mapping: 400
	ErrFailedPrecondition = werror.Wrap(nil).SetErrorf("FailedPrecondition") // HTTP Mapping: 400
	ErrOutOfRange         = werror.Wrap(nil).SetErrorf("OutOfRange")         // HTTP Mapping: 400
	ErrUnauthenticated    = werror.Wrap(nil).SetErrorf("Unauthenticated")    // HTTP Mapping: 401
	ErrPermissionDenied   = werror.Wrap(nil).SetErrorf("PermissionDenied")   // HTTP Mapping: 403
	ErrNotFound           = werror.Wrap(nil).SetErrorf("NotFound")           // HTTP Mapping: 404
	ErrAborted            = werror.Wrap(nil).SetErrorf("Aborted")            // HTTP Mapping: 409
	ErrAlreadyExists      = werror.Wrap(nil).SetErrorf("AlreadyExists")      // HTTP Mapping: 409
	ErrResourceExhausted  = werror.Wrap(nil).SetErrorf("ResourceExhausted")  // HTTP Mapping: 429
	ErrCancelled          = werror.Wrap(nil).SetErrorf("Cancelled")          // HTTP Mapping: 499
	ErrDataLoss           = werror.Wrap(nil).SetErrorf("DataLoss")           // HTTP Mapping: 500
	ErrUnknown            = werror.Wrap(nil).SetErrorf("Unknown")            // HTTP Mapping: 500
	ErrInternal           = werror.Wrap(nil).SetErrorf("Internal")           // HTTP Mapping: 500
	ErrNotImplemented     = werror.Wrap(nil).SetErrorf("NotImplemented")     // HTTP Mapping: 501
	ErrUnavailable        = werror.Wrap(nil).SetErrorf("Unavailable")        // HTTP Mapping: 503
	ErrDeadlineExceeded   = werror.Wrap(nil).SetErrorf("DeadlineExceeded")   // HTTP Mapping: 504
)
