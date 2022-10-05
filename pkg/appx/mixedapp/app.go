package mixedapp

import (
	"github.com/RussellLuo/appx"
	"github.com/RussellLuo/kun/pkg/appx/cronapp"
	"github.com/RussellLuo/kun/pkg/appx/cronapp2"
	"github.com/RussellLuo/kun/pkg/appx/httpapp"
)

type App struct {
	*appx.App
}

func New(name string, instance appx.Instance) *App {
	return &App{App: appx.New(name, instance)}
}

func (a *App) MountOn(parent, pattern string) *App {
	a.App.Use(httpapp.MountOn(parent, pattern))
	a.App.Require(parent)
	return a
}

func (a *App) ScheduledBy(scheduler, expression string) *App {
	a.App.Use(cronapp.ScheduledBy(a.App.Name, scheduler, expression))
	a.App.Require(scheduler)
	return a
}

func (a *App) ScheduledBy2(scheduler string) *App {
	a.App.Use(cronapp2.ScheduledBy(scheduler))
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
