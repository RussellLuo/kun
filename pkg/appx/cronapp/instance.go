package cronapp

import (
	"fmt"
)

// Scheduler represents a cron scheduler.
type Scheduler interface {
	Add(name, expr string, task func()) error
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

func GetCronScheduler(instance interface{}) (Scheduler, error) {
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

func GetCronJob(instance interface{}) (Job, error) {
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
