package cronapp

// Scheduler represents a cron scheduler.
type Scheduler interface {
	Add(name, expr string, task func()) error
}

// Job represents a cron job.
type Job interface {
	Task()
}
