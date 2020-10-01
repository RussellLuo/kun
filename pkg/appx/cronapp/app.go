package cronapp

import (
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

func (a *App) InitFunc(initFunc appx.InitFuncV2) *appx.App {
	init := initFunc
	if a.scheduler != "" {
		init = a.scheduledBy(initFunc)
	}

	a.App.InitFunc(init)
	return a.App // Return the wrapped *appx.App
}

func (a *App) scheduledBy(initFunc appx.InitFuncV2) appx.InitFuncV2 {
	return func(ctx appx.Context) error {
		if err := initFunc(ctx); err != nil {
			return err
		}

		job, ok := ctx.App.Value.(Job)
		if !ok {
			return fmt.Errorf("value %#v does not implement cronapp.Job", ctx.App.Value)
		}

		schedulerValue := ctx.Required[a.scheduler].Value
		scheduler, ok := schedulerValue.(Scheduler)
		if !ok {
			return fmt.Errorf("value %#v does not implement cronapp.Scheduler", schedulerValue)
		}

		scheduler.Add(a.Name, a.expression, job.Task) // nolint:errcheck
		return nil
	}
}
