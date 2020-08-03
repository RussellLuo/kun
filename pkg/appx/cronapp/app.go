package cronapp

import (
	"context"
	"fmt"

	"github.com/RussellLuo/appx"
)

type App struct {
	*appx.App

	scheduler  string
	expression string
}

func New(name string) *App {
	return &App{App: appx.New(name)}
}

func (a *App) ScheduledBy(scheduler string) *App {
	a.scheduler = scheduler
	a.App.Require(scheduler)
	return a
}

func (a *App) Expression(expression string) *App {
	a.expression = expression
	return a
}

func (a *App) Require(names ...string) *App {
	a.App.Require(names...)
	return a
}

func (a *App) Init(initFunc appx.InitFunc) *appx.App {
	init := initFunc
	if a.scheduler != "" {
		init = a.scheduledBy(initFunc)
	}

	a.App.Init2(init)
	return a.App // Return the wrapped *appx.App
}

func (a *App) scheduledBy(initFunc appx.InitFunc) appx.InitFunc {
	return func(ctx context.Context, lc appx.Lifecycle, apps map[string]*appx.App) (appx.Value, appx.CleanFunc, error) {
		value, clean, err := initFunc(ctx, lc, apps)
		if err != nil {
			return nil, nil, err
		}

		job, ok := value.(Job)
		if !ok {
			return nil, nil, fmt.Errorf("value %#v does not implement cronapp.Job", value)
		}

		schedulerValue := apps[a.scheduler].Value
		scheduler, ok := schedulerValue.(Scheduler)
		if !ok {
			return nil, nil, fmt.Errorf("value %#v does not implement cronapp.Scheduler", schedulerValue)
		}

		scheduler.Add(a.Name, a.expression, job.Task) // nolint:errcheck
		return value, clean, nil
	}
}
