package tickdoer

import (
	"sync"
	"time"
)

type TickDoer struct {
	t *time.Ticker
	f func()

	exitC     chan struct{}
	waitGroup sync.WaitGroup
}

func TickFunc(d time.Duration, f func()) *TickDoer {
	td := &TickDoer{
		t: time.NewTicker(d),
		f: f,
	}
	td.start()
	return td
}

func (td *TickDoer) start() {
	td.waitGroup.Add(1)
	go func() {
		defer td.waitGroup.Done()

		for {
			select {
			case <-td.t.C:
				td.f()
			case <-td.exitC:
				return
			}
		}
	}()
}

func (td *TickDoer) Stop() {
	td.t.Stop()

	close(td.exitC)
	td.waitGroup.Wait()
}
