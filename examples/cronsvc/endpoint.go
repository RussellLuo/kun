// Code generated by kun; DO NOT EDIT.
// github.com/RussellLuo/kun

package cronsvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type SendEmailResponse struct {
	Err error `json:"-"`
}

func (r *SendEmailResponse) Body() interface{} { return r }

// Failed implements endpoint.Failer.
func (r *SendEmailResponse) Failed() error { return r.Err }

// MakeEndpointOfSendEmail creates the endpoint for s.SendEmail.
func MakeEndpointOfSendEmail(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.SendEmail(
			ctx,
		)
		return &SendEmailResponse{
			Err: err,
		}, nil
	}
}
