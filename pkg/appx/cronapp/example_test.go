package cronapp_test

import (
	"context"
	"fmt"
	"time"

	"github.com/RussellLuo/appx"
	"github.com/RussellLuo/micron/cron"
	nillocker "github.com/RussellLuo/micron/locker/nil"

	"github.com/RussellLuo/kok/pkg/appx/cronapp"
)

type task func()

func (t task) Task() {
	t()
}

func Example() {
	// Typically located in `func init()` of package hi.
	appx.MustRegister(cronapp.New("hi").
		ScheduledBy("greeter").Expression("*/1 * * * * * *").
		Init(func(ctx context.Context, lc appx.Lifecycle, apps map[string]*appx.App) (appx.Value, appx.CleanFunc, error) {
			return task(func() {
				fmt.Println("Saying hi to you")
			}), nil, nil
		}))

	// Typically located in `func init()` of package bye.
	appx.MustRegister(cronapp.New("bye").
		ScheduledBy("greeter").Expression("*/2 * * * * * *").
		Init(func(ctx context.Context, lc appx.Lifecycle, apps map[string]*appx.App) (appx.Value, appx.CleanFunc, error) {
			return task(func() {
				fmt.Println("Saying bye to you")
			}), nil, nil
		}))

	// Typically located in `func init()` of package greeter.
	appx.MustRegister(cronapp.New("greeter").
		Init(func(ctx context.Context, lc appx.Lifecycle, apps map[string]*appx.App) (appx.Value, appx.CleanFunc, error) {
			c := cron.New(nillocker.New(), nil)
			lc.Append(appx.Hook{
				OnStart: func(context.Context) error {
					fmt.Println("Starting CRON scheduler")
					c.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					fmt.Println("Stopping CRON scheduler")
					c.Stop()
					return nil
				},
			})
			return c, nil, nil
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

	// Wait for 2.01s to give the two jobs enough time to run at least once.
	time.Sleep(2010 * time.Millisecond)

	// Stop the greeter.
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	appx.Stop(stopCtx)

	// Output:
	// Starting CRON scheduler
	// Saying hi to you
	// Saying bye to you
	// Saying hi to you
	// Stopping CRON scheduler
}
