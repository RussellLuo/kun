package httpapp

import (
	"errors"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
)

func TestGetService(t *testing.T) {
	svc := struct{}{}
	cases := []struct {
		in      interface{}
		wantSvc interface{}
		wantErr error
	}{
		{
			in:      nil,
			wantSvc: nil,
			wantErr: errors.New("value <nil> cannot be converted to *httpapp.Value"),
		},
		{
			in:      (*Value)(nil),
			wantSvc: nil,
			wantErr: errors.New("value (*httpapp.Value)(nil) holds no service"),
		},
		{
			in:      &Value{},
			wantSvc: nil,
			wantErr: errors.New("value &httpapp.Value{Service:interface {}(nil), Router:chi.Router(nil)} holds no service"),
		},
		{
			in:      &Value{Service: svc},
			wantSvc: svc,
			wantErr: nil,
		},
	}
	for _, c := range cases {
		svc, err := GetService(c.in)
		if svc != c.wantSvc {
			t.Fatalf("Service: got (%#v), want (%#v)", svc, c.wantSvc)
		}
		if !reflect.DeepEqual(err, c.wantErr) {
			t.Fatalf("Error: got (%#v), want (%#v)", err, c.wantErr)
		}
	}
}

func TestGetRouter(t *testing.T) {
	router := chi.NewRouter()
	cases := []struct {
		in         interface{}
		wantRouter interface{}
		wantErr    error
	}{
		{
			in:         nil,
			wantRouter: nil,
			wantErr:    errors.New("value <nil> cannot be converted to *httpapp.Value"),
		},
		{
			in:         (*Value)(nil),
			wantRouter: nil,
			wantErr:    errors.New("value (*httpapp.Value)(nil) is not routable"),
		},
		{
			in:         &Value{},
			wantRouter: nil,
			wantErr:    errors.New("value &httpapp.Value{Service:interface {}(nil), Router:chi.Router(nil)} is not routable"),
		},
		{
			in:         &Value{Router: router},
			wantRouter: router,
			wantErr:    nil,
		},
	}
	for _, c := range cases {
		router, err := GetRouter(c.in)
		if router != c.wantRouter {
			t.Fatalf("Router: got (%#v), want (%#v)", router, c.wantRouter)
		}
		if !reflect.DeepEqual(err, c.wantErr) {
			t.Fatalf("Error: got (%#v), want (%#v)", err, c.wantErr)
		}
	}
}
