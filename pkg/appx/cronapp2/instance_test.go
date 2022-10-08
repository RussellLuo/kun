package cronapp2

import (
	"errors"
	"reflect"
	"testing"

	"github.com/RussellLuo/kun/pkg/appx/cronapp"
	"github.com/RussellLuo/micron"
)

type scheduler struct{}

func (s *scheduler) Add(name, expr string, task func()) error {
	return nil
}

func (s *scheduler) AddJob(job ...micron.Job) error {
	return nil
}

type cronScheduler struct {
	scheduler cronapp.Scheduler
}

func (cs *cronScheduler) Scheduler() cronapp.Scheduler {
	return cs.scheduler
}

type cronJobs struct {
	jobs []micron.Job
}

func (cj *cronJobs) Jobs() []micron.Job {
	return cj.jobs
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
			wantErr:       errors.New("instance <nil> does not implement cronapp.CronScheduler"),
		},
		{
			in:            &cronScheduler{scheduler: nil},
			wantScheduler: nil,
			wantErr:       errors.New("method Scheduler() of instance &cronapp2.cronScheduler{scheduler:cronapp.Scheduler(nil)} returns nil"),
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

func TestGetCronJobs(t *testing.T) {
	nopJobs := []micron.Job{
		{
			Task: func() {},
		},
	}

	cases := []struct {
		in       interface{}
		wantJobs []micron.Job
		wantErr  error
	}{
		{
			in:       nil,
			wantJobs: nil,
			wantErr:  errors.New("instance <nil> does not implement cronapp2.CronJobs"),
		},
		{
			in:       &cronJobs{},
			wantJobs: nil,
			wantErr:  errors.New("method Jobs() of value &cronapp2.cronJobs{jobs:[]micron.Job(nil)} returns nil"),
		},
		{
			in:       &cronJobs{jobs: nopJobs},
			wantJobs: nopJobs,
			wantErr:  nil,
		},
	}
	for _, c := range cases {
		jobs, err := getCronJobs(c.in)
		if !reflect.DeepEqual(jobs, c.wantJobs) {
			t.Fatalf("Jobs: got (%#v), want (%#v)", jobs, c.wantJobs)
		}
		if !reflect.DeepEqual(err, c.wantErr) {
			t.Fatalf("Error: got (%#v), want (%#v)", err, c.wantErr)
		}
	}
}
