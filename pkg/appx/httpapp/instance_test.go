package httpapp

import (
	"errors"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
)

type instance struct {
	router chi.Router
}

func (i *instance) Router() chi.Router {
	return i.router
}

func TestGetChiRouter(t *testing.T) {
	router := chi.NewRouter()

	cases := []struct {
		in         interface{}
		wantRouter interface{}
		wantErr    error
	}{
		{
			in:         nil,
			wantRouter: nil,
			wantErr:    errors.New("instance <nil> does not implement httpapp.ChiRouter"),
		},
		{
			in:         &instance{router: nil},
			wantRouter: nil,
			wantErr:    errors.New("method Router() of instance &httpapp.instance{router:chi.Router(nil)} returns nil"),
		},
		{
			in:         &instance{router: router},
			wantRouter: router,
			wantErr:    nil,
		},
	}
	for _, c := range cases {
		router, err := getChiRouter(c.in)
		if router != c.wantRouter {
			t.Fatalf("Router: got (%#v), want (%#v)", router, c.wantRouter)
		}
		if !reflect.DeepEqual(err, c.wantErr) {
			t.Fatalf("Error: got (%#v), want (%#v)", err, c.wantErr)
		}
	}
}
