package httpapp_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RussellLuo/appx"
	"github.com/go-chi/chi"

	"github.com/RussellLuo/kok/pkg/appx/httpapp"
)

func Example() {
	// Typically located in `func init()` of package hi.
	appx.MustRegister(httpapp.New("hi").
		MountOn("greeter").Pattern("/hi").
		Init(func(ctx context.Context, lc appx.Lifecycle, apps map[string]*appx.App) (appx.Value, appx.CleanFunc, error) {
			r := chi.NewRouter()
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("Got a request for /hi")
			})
			return &httpapp.Value{
				Router: r,
			}, nil, nil
		}))

	// Typically located in `func init()` of package bye.
	appx.MustRegister(httpapp.New("bye").
		MountOn("greeter").Pattern("/bye").
		Init(func(ctx context.Context, lc appx.Lifecycle, apps map[string]*appx.App) (appx.Value, appx.CleanFunc, error) {
			r := chi.NewRouter()
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("Got a request for /bye")
			})
			return &httpapp.Value{
				Router: r,
			}, nil, nil
		}))

	// Typically located in `func init()` of package greeter.
	appx.MustRegister(httpapp.New("greeter").
		Init(func(ctx context.Context, lc appx.Lifecycle, apps map[string]*appx.App) (appx.Value, appx.CleanFunc, error) {
			r := chi.NewRouter()
			server := &http.Server{
				Addr:    ":8080",
				Handler: r,
			}
			lc.Append(appx.Hook{
				OnStart: func(context.Context) error {
					fmt.Println("Starting HTTP server")
					go server.ListenAndServe() // nolint:errcheck
					return nil
				},
				OnStop: func(ctx context.Context) error {
					fmt.Println("Stopping HTTP server")
					return server.Shutdown(ctx)
				},
			})
			return &httpapp.Value{
				Router: r,
			}, nil, nil
		}))

	// Typically located in `func main()` of package main.
	appx.ErrorHandler(func(err error) {
		fmt.Printf("err: %v\n", err)
	})

	// Installs the applications.
	if err := appx.Install(context.Background()); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer appx.Uninstall()

	// Start the greeter.
	startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := appx.Start(startCtx); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	// Make two HTTP requests to demonstrate that our server is running.
	http.Get("http://localhost:8080/hi")  // nolint:errcheck
	http.Get("http://localhost:8080/bye") // nolint:errcheck

	// Stop the greeter.
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	appx.Stop(stopCtx)

	// Output:
	// Starting HTTP server
	// Got a request for /hi
	// Got a request for /bye
	// Stopping HTTP server
}
