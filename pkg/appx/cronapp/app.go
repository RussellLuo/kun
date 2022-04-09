package cronapp

import (
	"github.com/RussellLuo/appx"
)

type App struct {
	*appx.App
}

func New(name string, instance appx.Instance) *App {
	return &App{App: appx.New(name, instance)}
}

func (a *App) ScheduledBy(scheduler, expression string) *App {
	m := ScheduledBy(a.App.Name, scheduler, expression)
	a.App.Instance = appx.Standardize(m(a.App.Instance))

	a.App.Require(scheduler)
	return a
}

func (a *App) Require(names ...string) *App {
	a.App.Require(names...)
	return a
}
