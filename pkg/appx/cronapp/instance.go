package cronapp

import (
	"fmt"

	"github.com/RussellLuo/appx"
	"github.com/RussellLuo/micron/cron"
)

func ScheduledBy(name, scheduler, expression string) func(appx.Standard) appx.Standard {
	return func(next appx.Standard) appx.Standard {
		return middleware{
			Standard:   next,
			name:       name,
			scheduler:  scheduler,
			expression: expression,
		}
	}
}

type middleware struct {
	appx.Standard
	name       string
	scheduler  string
	expression string
}

func (m middleware) Init(ctx appx.Context) error {
	if err := m.Standard.Init(ctx); err != nil {
		return err
	}

	scheduler, err := getCronScheduler(ctx.MustLoad(m.scheduler))
	if err != nil {
		return err
	}

	job, err := getCronJob(m.Standard.Instance())
	if err != nil {
		return err
	}

	_ = scheduler.Add(m.name, m.expression, job.Task)
	return nil
}

// Scheduler represents a cron scheduler.
type Scheduler interface {
	Add(name, expr string, task func()) error
	AddJob(job ...cron.Job) error
}

// Job represents a cron job.
type Job interface {
	Task()
}

// CronScheduler is the interface that a scheduler application must implement.
type CronScheduler interface {
	Scheduler() Scheduler
}

// CronJob is the interface that a job application must implement.
type CronJob interface {
	Job() Job
}

func getCronScheduler(instance interface{}) (Scheduler, error) {
	r, ok := instance.(CronScheduler)
	if !ok {
		return nil, fmt.Errorf("instance %#v does not implement httpapp.CronScheduler", instance)
	}

	result := r.Scheduler()
	if result == nil {
		return nil, fmt.Errorf("method Scheduler() of instance %#v returns nil", instance)
	}

	return result, nil
}

func getCronJob(instance interface{}) (Job, error) {
	r, ok := instance.(CronJob)
	if !ok {
		return nil, fmt.Errorf("instance %#v does not implement httpapp.CronJob", instance)
	}

	result := r.Job()
	if result == nil {
		return nil, fmt.Errorf("method Job() of value %#v returns nil", instance)
	}

	return result, nil
}
