package cronapp2

import (
	"fmt"

	"github.com/RussellLuo/appx"
	"github.com/RussellLuo/kun/pkg/appx/cronapp"
	"github.com/RussellLuo/micron"
)

func ScheduledBy(scheduler string) func(appx.Standard) appx.Standard {
	return func(next appx.Standard) appx.Standard {
		return middleware{
			Standard:  next,
			scheduler: scheduler,
		}
	}
}

type middleware struct {
	appx.Standard
	scheduler string
}

func (m middleware) Init(ctx appx.Context) error {
	if err := m.Standard.Init(ctx); err != nil {
		return err
	}

	scheduler, err := getCronScheduler(ctx.MustLoad(m.scheduler))
	if err != nil {
		return err
	}

	jobs, err := getCronJobs(m.Standard.Instance())
	if err != nil {
		return err
	}

	return scheduler.AddJob(jobs...)
}

// CronScheduler is the interface that a scheduler application must implement.
type CronScheduler interface {
	Scheduler() cronapp.Scheduler
}

// CronJobs is the interface that a job application must implement.
type CronJobs interface {
	Jobs() []micron.Job
}

func getCronScheduler(instance interface{}) (cronapp.Scheduler, error) {
	r, ok := instance.(CronScheduler)
	if !ok {
		return nil, fmt.Errorf("instance %#v does not implement cronapp.CronScheduler", instance)
	}

	result := r.Scheduler()
	if result == nil {
		return nil, fmt.Errorf("method Scheduler() of instance %#v returns nil", instance)
	}

	return result, nil
}

func getCronJobs(instance interface{}) ([]micron.Job, error) {
	r, ok := instance.(CronJobs)
	if !ok {
		return nil, fmt.Errorf("instance %#v does not implement cronapp2.CronJobs", instance)
	}

	result := r.Jobs()
	if result == nil {
		return nil, fmt.Errorf("method Jobs() of value %#v returns nil", instance)
	}

	return result, nil
}
