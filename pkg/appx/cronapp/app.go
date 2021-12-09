package cronapp

import (
	"github.com/RussellLuo/appx"
	"github.com/RussellLuo/kun/pkg/appx/wrapper"
)

type App struct {
	*appx.App

	wrapper *wrapper.InitWrapper

	scheduler  string
	expression string
}

func New(name string, instance appx.Initializer) *App {
	w := wrapper.New(instance)
	return &App{
		App:     appx.New(name, w),
		wrapper: w,
	}
}

func (a *App) ScheduledBy(scheduler string, expression string) *App {
	a.scheduler = scheduler
	a.App.Require(scheduler)

	a.expression = expression
	a.wrapper.AfterInitFunc(a.afterInit)

	return a
}

func (a *App) Require(names ...string) *App {
	a.App.Require(names...)
	return a
}

func (a *App) afterInit(ctx appx.Context) error {
	scheduler, err := GetCronScheduler(ctx.MustLoad(a.scheduler))
	if err != nil {
		return err
	}

	job, err := GetCronJob(a.wrapper.Instance())
	if err != nil {
		return err
	}

	scheduler.Add(a.Name, a.expression, job.Task) // nolint:errcheck
	return nil
}
