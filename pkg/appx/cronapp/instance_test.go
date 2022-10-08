package cronapp

import (
	"errors"
	"reflect"
	"testing"

	"github.com/RussellLuo/micron"
)

type scheduler struct{}

func (s *scheduler) Add(name, expr string, task func()) error {
	return nil
}

func (s *scheduler) AddJob(job ...micron.Job) error {
	return nil
}

type job struct {
	task func()
}

func (j *job) Task() {
	j.task()
}

type cronScheduler struct {
	scheduler Scheduler
}

func (cs *cronScheduler) Scheduler() Scheduler {
	return cs.scheduler
}

type cronJob struct {
	job Job
}

func (cj *cronJob) Job() Job {
	return cj.job
}

func TestGetCronScheduler(t *testing.T) {
	nopScheduler := &scheduler{}

	cases := []struct {
		in            interface{}
		wantScheduler interface{}
		wantErr       error
	}{
		{
			in:            nil,
			wantScheduler: nil,
			wantErr:       errors.New("instance <nil> does not implement httpapp.CronScheduler"),
		},
		{
			in:            &cronScheduler{scheduler: nil},
			wantScheduler: nil,
			wantErr:       errors.New("method Scheduler() of instance &cronapp.cronScheduler{scheduler:cronapp.Scheduler(nil)} returns nil"),
		},
		{
			in:            &cronScheduler{scheduler: nopScheduler},
			wantScheduler: nopScheduler,
			wantErr:       nil,
		},
	}
	for _, c := range cases {
		router, err := getCronScheduler(c.in)
		if router != c.wantScheduler {
			t.Fatalf("Scheduler: got (%#v), want (%#v)", router, c.wantScheduler)
		}
		if !reflect.DeepEqual(err, c.wantErr) {
			t.Fatalf("Error: got (%#v), want (%#v)", err, c.wantErr)
		}
	}
}

func TestGetCronJob(t *testing.T) {
	nopJob := &job{task: func() {}}

	cases := []struct {
		in      interface{}
		wantJob interface{}
		wantErr error
	}{
		{
			in:      nil,
			wantJob: nil,
			wantErr: errors.New("instance <nil> does not implement httpapp.CronJob"),
		},
		{
			in:      &cronJob{job: nil},
			wantJob: nil,
			wantErr: errors.New("method Job() of value &cronapp.cronJob{job:cronapp.Job(nil)} returns nil"),
		},
		{
			in:      &cronJob{job: nopJob},
			wantJob: nopJob,
			wantErr: nil,
		},
	}
	for _, c := range cases {
		router, err := getCronJob(c.in)
		if router != c.wantJob {
			t.Fatalf("Job: got (%#v), want (%#v)", router, c.wantJob)
		}
		if !reflect.DeepEqual(err, c.wantErr) {
			t.Fatalf("Error: got (%#v), want (%#v)", err, c.wantErr)
		}
	}
}
