package httpapp_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RussellLuo/appx"
	"github.com/go-chi/chi"

	"github.com/RussellLuo/kun/pkg/appx/httpapp"
)

type Hi struct {
	router chi.Router
}

func (h *Hi) Router() chi.Router {
	return h.router
}

func (h *Hi) Init(appx.Context) error {
	h.router = chi.NewRouter()
	h.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got a request for /hi")
	})
	return nil
}

type Bye struct {
	router chi.Router
}

func (b *Bye) Router() chi.Router {
	return b.router
}

func (b *Bye) Init(appx.Context) error {
	b.router = chi.NewRouter()
	b.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got a request for /bye")
	})
	return nil
}

type Greeter struct {
	router chi.Router
	server *http.Server
}

func (g *Greeter) Router() chi.Router {
	return g.router
}

func (g *Greeter) Init(appx.Context) error {
	g.router = chi.NewRouter()
	g.server = &http.Server{
		Addr:    ":8080",
		Handler: g.router,
	}
	return nil
}

func (g *Greeter) Start(context.Context) error {
	fmt.Println("Starting HTTP server")
	go g.server.ListenAndServe() // nolint:errcheck
	return nil
}

func (g *Greeter) Stop(ctx context.Context) error {
	fmt.Println("Stopping HTTP server")
	return g.server.Shutdown(ctx)
}

func Example() {
	r := appx.NewRegistry()

	// Typically located in `func init()` of package hi.
	r.MustRegister(httpapp.New("hi", new(Hi)).MountOn("greeter", "/hi").App)

	// Typically located in `func init()` of package bye.
	r.MustRegister(httpapp.New("bye", new(Bye)).MountOn("greeter", "/bye").App)

	// Typically located in `func init()` of package greeter.
	r.MustRegister(httpapp.New("greeter", new(Greeter)).App)

	// Typically located in `func main()` of package main.
	r.SetOptions(&appx.Options{
		ErrorHandler: func(err error) {
			fmt.Printf("err: %v\n", err)
		},
	})

	// Installs the applications.
	if err := r.Install(context.Background()); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer r.Uninstall()

	// Start the greeter.
	startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.Start(startCtx); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	// Make two HTTP requests to demonstrate that our server is running.
	http.Get("http://localhost:8080/hi")  // nolint:errcheck
	http.Get("http://localhost:8080/bye") // nolint:errcheck

	// Stop the greeter.
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r.Stop(stopCtx)

	// Output:
	// Starting HTTP server
	// Got a request for /hi
	// Got a request for /bye
	// Stopping HTTP server
}
