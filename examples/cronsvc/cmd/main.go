package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RussellLuo/kun/examples/cronsvc"
	"github.com/RussellLuo/micron"
)

func main() {
	c := micron.New(
		micron.NewSemaphoreLocker(),
		&micron.Options{
			Timezone: "Asia/Shanghai",
			LockTTL:  2 * time.Second, // Assume the maximal clock error is 2s.
			ErrHandler: func(err error) {
				log.Printf("err: %v", err)
			},
		},
	)

	jobs := cronsvc.NewCronJobs(&cronsvc.Handler{})
	if err := c.AddJob(jobs...); err != nil {
		log.Printf("err: %v\n", err)
	}

	c.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	c.Stop()
}
