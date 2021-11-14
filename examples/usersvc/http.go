// Code generated by kok; DO NOT EDIT.
// github.com/RussellLuo/kok

package usersvc

import (
	"context"
	"net/http"

	"github.com/RussellLuo/kok/pkg/httpcodec"
	"github.com/RussellLuo/kok/pkg/httpoption"
	"github.com/RussellLuo/kok/pkg/oasv2"
	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHTTPRouter(svc Service, codecs httpcodec.Codecs, opts ...httpoption.Option) chi.Router {
	r := chi.NewRouter()
	options := httpoption.NewOptions(opts...)

	r.Method("GET", "/api", oasv2.Handler(OASv2APIDoc, options.ResponseSchema()))

	var codec httpcodec.Codec
	var validator httpoption.Validator
	var kitOptions []kithttp.ServerOption

	codec = codecs.EncodeDecoder("CreateUser")
	validator = options.RequestValidator("CreateUser")
	r.Method(
		"POST", "/users",
		kithttp.NewServer(
			MakeEndpointOfCreateUser(svc),
			decodeCreateUserRequest(codec, validator),
			httpcodec.MakeResponseEncoder(codec, 200),
			append(kitOptions,
				kithttp.ServerErrorEncoder(httpcodec.MakeErrorEncoder(codec)),
			)...,
		),
	)

	return r
}

func NewHTTPRouterWithOAS(svc Service, codecs httpcodec.Codecs, schema oasv2.Schema) chi.Router {
	return NewHTTPRouter(svc, codecs, httpoption.ResponseSchema(schema))
}

func decodeCreateUserRequest(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var _req CreateUserRequest

		user := map[string][]string{
			"query.name":             r.URL.Query()["name"],
			"query.age":              r.URL.Query()["age"],
			"header.X-Forwarded-For": r.Header.Values("X-Forwarded-For"),
			"request.RemoteAddr":     []string{r.RemoteAddr},
		}
		if err := codec.DecodeRequestParams("user", user, &_req.User); err != nil {
			return nil, err
		}

		if err := validator.Validate(&_req); err != nil {
			return nil, err
		}

		return &_req, nil
	}
}
