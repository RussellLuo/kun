package cronapp

import (
	"fmt"
)

// Scheduler represents a cron scheduler.
type Scheduler interface {
	Add(name, expr string, task func()) error
}

func GetScheduler(value interface{}) (Scheduler, error) {
	s, ok := value.(Scheduler)
	if !ok {
		return nil, fmt.Errorf("value %#v does not implement cronapp.Scheduler", value)
	}
	return s, nil
}

// Job represents a cron job.
type Job interface {
	Task()
}

func GetJob(value interface{}) (Job, error) {
	j, ok := value.(Job)
	if !ok {
		return nil, fmt.Errorf("value %#v does not implement cronapp.Job", value)
	}
	return j, nil
}

// CronScheduler is the interface that a scheduler application must implement.
type CronScheduler interface {
	Scheduler() Scheduler
}

// CronJob is the interface that a job application must implement.
type CronJob interface {
	Job() Job
}

func GetCronScheduler(value interface{}) (Scheduler, error) {
	r, ok := value.(CronScheduler)
	if !ok {
		return nil, fmt.Errorf("value %#v does not implement httpapp.CronScheduler", value)
	}

	result := r.Scheduler()
	if result == nil {
		return nil, fmt.Errorf("method Scheduler() of value %#v returns nil", value)
	}

	return result, nil
}

func GetCronJob(value interface{}) (Job, error) {
	r, ok := value.(CronJob)
	if !ok {
		return nil, fmt.Errorf("value %#v does not implement httpapp.CronJob", value)
	}

	result := r.Job()
	if result == nil {
		return nil, fmt.Errorf("method Job() of value %#v returns nil", value)
	}

	return result, nil
}
