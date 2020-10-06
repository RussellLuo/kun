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
