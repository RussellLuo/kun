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
	words := make(chan string, 3)
	sleepAndPrintTimes := func() {
		// Wait for 2.02s to execute the jobs 3 times.
		time.Sleep(2020 * time.Millisecond)

		// Calculate and print the execution times.
		times := map[string]int{
			"hi":  0,
			"bye": 0,
		}
		for i := 0; i < cap(words); i++ {
			times[<-words]++
		}
		fmt.Printf("Saying hi %d times\n", times["hi"])
		fmt.Printf("Saying bye %d time\n", times["bye"])
	}

	// Typically located in `func init()` of package hi.
	appx.MustRegister(
		cronapp.New("hi").
			ScheduledBy("greeter").Expression("*/1 * * * * * *").
			InitFunc(func(ctx appx.Context) error {
				ctx.App.Value = task(func() {
					words <- "hi"
				})
				return nil
			}),
	)

	// Typically located in `func init()` of package bye.
	appx.MustRegister(
		cronapp.New("bye").
			ScheduledBy("greeter").Expression("*/2 * * * * * *").
			InitFunc(func(ctx appx.Context) error {
				ctx.App.Value = task(func() {
					words <- "bye"
				})
				return nil
			}),
	)

	// Typically located in `func init()` of package greeter.
	appx.MustRegister(
		cronapp.New("greeter").
			InitFunc(func(ctx appx.Context) error {
				c := cron.New(nillocker.New(), nil)
				ctx.Lifecycle.Append(appx.Hook{
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
				ctx.App.Value = c
				return nil
			}),
	)

	// Typically located in `func main()` of package main.
	appx.SetConfig(appx.Config{
		ErrorHandler: func(err error) {
			fmt.Printf("err: %v\n", err)
		},
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

	sleepAndPrintTimes()

	// Stop the greeter.
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	appx.Stop(stopCtx)

	// Output:
	// Starting CRON scheduler
	// Saying hi 2 times
	// Saying bye 1 time
	// Stopping CRON scheduler
}
