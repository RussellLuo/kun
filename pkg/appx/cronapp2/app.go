package cronapp2

import (
	"github.com/RussellLuo/appx"
)

type App struct {
	*appx.App
}

func New(name string, instance appx.Instance) *App {
	return &App{App: appx.New(name, instance)}
}

func (a *App) ScheduledBy(scheduler string) *App {
	a.App.Use(ScheduledBy(scheduler))
	a.App.Require(scheduler)
	return a
}

func (a *App) Use(middlewares ...func(appx.Standard) appx.Standard) *App {
	a.App.Use(middlewares...)
	return a
}

func (a *App) Require(names ...string) *App {
	a.App.Require(names...)
	return a
}
