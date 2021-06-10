package cronapp

import (
	"github.com/RussellLuo/appx"
	"github.com/RussellLuo/kok/pkg/appx/wrapper"
)

type App struct {
	*appx.App

	wrapper *wrapper.InitWrapper

	scheduler  string
	expression string
}

func New(name string) *App {
	return &App{App: appx.New(name)}
}

func NewV2(name string, instance appx.Initializer) *App {
	w := wrapper.New(instance)
	return &App{
		App:     appx.NewV2(name, w),
		wrapper: w,
	}
}

func (a *App) ScheduledBy(scheduler string, expression ...string) *App {
	a.scheduler = scheduler
	a.App.Require(scheduler)

	if len(expression) == 1 {
		a.expression = expression[0]
		a.wrapper.AfterInitFunc(a.afterInit)
	}

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

		scheduler, err := GetScheduler(ctx.Required[a.scheduler].Value)
		if err != nil {
			return err
		}

		job, err := GetJob(ctx.App.Value)
		if err != nil {
			return err
		}

		scheduler.Add(a.Name, a.expression, job.Task) // nolint:errcheck
		return nil
	}
}

func (a *App) afterInit(ctx appx.Context) error {
	var scheduler Scheduler
	var err error

	instance := ctx.MustLoad(a.scheduler)
	if inst, ok := instance.(CronScheduler); ok {
		scheduler, err = GetCronScheduler(inst)
		if err != nil {
			return err
		}
	} else {
		// The parent app is an old-style one.
		// TODO: remove this snippet when all apps are new-style ones.
		scheduler, err = GetScheduler(ctx.Required[a.scheduler].Value)
		if err != nil {
			return err
		}
	}

	job, err := GetCronJob(a.wrapper.Instance())
	if err != nil {
		return err
	}

	scheduler.Add(a.Name, a.expression, job.Task) // nolint:errcheck
	return nil
}
