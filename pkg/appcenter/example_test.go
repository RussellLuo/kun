package appcenter_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RussellLuo/appx"
	"github.com/go-chi/chi"

	"github.com/RussellLuo/kok/pkg/appcenter"
)

func Example() {

	// Typically located in `func init()` of package b.
	appx.MustRegister(appcenter.New("hi").
		MountOn("main").Pattern("/v1/hi").
		Init(func(ctx context.Context, apps map[string]*appx.App) (appx.Value, appx.CleanFunc, error) {
			r := chi.NewRouter()
			r.Get("/h", nil)
			r.Get("/i", nil)
			return &appcenter.Value{
				Router: r,
			}, nil, nil
		}))

	// Typically located in `func init()` of package b.
	appx.MustRegister(appcenter.New("bye").
		MountOn("main").Pattern("/v1/bye").
		Init(func(ctx context.Context, apps map[string]*appx.App) (appx.Value, appx.CleanFunc, error) {
			r := chi.NewRouter()
			r.Get("/b", nil)
			r.Get("/y", nil)
			r.Get("/e", nil)
			return &appcenter.Value{
				Router: r,
			}, nil, nil
		}))

	// Typically located in `func main()` of package main.
	r := chi.NewRouter()
	appx.MustRegister(appx.New("main").
		Init(func(ctx context.Context, apps map[string]*appx.App) (appx.Value, appx.CleanFunc, error) {
			return &appcenter.Value{
				Router: r,
			}, nil, nil
		}))

	if err := appx.Install(context.Background()); err != nil {
		fmt.Printf("err: %v\n", err)
	}

	// Walk the routes in router r.
	walkFunc := func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		fmt.Printf("method: %s, route: %s\n", method, route)
		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("err: %v\n", err)
	}

	if err := appx.Uninstall(); err != nil {
		fmt.Printf("err: %v\n", err)
	}

	// Output:
	// method: GET, route: /v1/bye/b
	// method: GET, route: /v1/bye/e
	// method: GET, route: /v1/bye/y
	// method: GET, route: /v1/hi/h
	// method: GET, route: /v1/hi/i
}
