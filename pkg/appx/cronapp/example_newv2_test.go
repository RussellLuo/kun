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

type Hi struct {
	words chan<- string
	job   cronapp.Job
}

func newHi(words chan<- string) *Hi {
	return &Hi{words: words}
}

func (h *Hi) Job() cronapp.Job {
	return h.job
}

func (h *Hi) Init(ctx appx.Context) error {
	h.job = task(func() {
		h.words <- "hi"
	})
	return nil
}

type Bye struct {
	words chan<- string
	job   cronapp.Job
}

func newBye(words chan<- string) *Bye {
	return &Bye{words: words}
}

func (b *Bye) Job() cronapp.Job {
	return b.job
}

func (b *Bye) Init(ctx appx.Context) error {
	b.job = task(func() {
		b.words <- "bye"
	})
	return nil
}

type Greeter struct {
	c *cron.Cron
}

func (g *Greeter) Scheduler() cronapp.Scheduler {
	return g.c
}

func (g *Greeter) Init(ctx appx.Context) error {
	g.c = cron.New(nillocker.New(), nil)
	return nil
}

func (g *Greeter) Start(ctx context.Context) error {
	fmt.Println("Starting CRON scheduler")
	g.c.Start()
	return nil
}

func (g *Greeter) Stop(ctx context.Context) error {
	fmt.Println("Stopping CRON scheduler")
	g.c.Stop()
	return nil
}

func Example_newV2() {
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

	r := appx.NewRegistry()

	// Typically located in `func init()` of package hi.
	r.MustRegister(
		cronapp.NewV2("hi", newHi(words)).
			ScheduledBy("greeter", "*/1 * * * * * *").App,
	)

	// Typically located in `func init()` of package bye.
	r.MustRegister(
		cronapp.NewV2("bye", newBye(words)).
			ScheduledBy("greeter", "*/2 * * * * * *").App,
	)

	// Typically located in `func init()` of package greeter.
	r.MustRegister(
		cronapp.NewV2("greeter", new(Greeter)).App,
	)

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

	sleepAndPrintTimes()

	// Stop the greeter.
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r.Stop(stopCtx)

	// Output:
	// Starting CRON scheduler
	// Saying hi 2 times
	// Saying bye 1 time
	// Stopping CRON scheduler
}
