package mixedapp

import (
	"github.com/RussellLuo/appx"
	"github.com/RussellLuo/kun/pkg/appx/cronapp"
	"github.com/RussellLuo/kun/pkg/appx/httpapp"
)

type App struct {
	*appx.App
}

func New(name string, instance appx.Instance) *App {
	return &App{App: appx.New(name, instance)}
}

func (a *App) MountOn(parent, pattern string) *App {
	m := httpapp.MountOn(parent, pattern)
	a.App.Instance = appx.Standardize(m(a.App.Instance))

	a.App.Require(parent)
	return a
}

func (a *App) ScheduledBy(scheduler, expression string) *App {
	m := cronapp.ScheduledBy(a.App.Name, scheduler, expression)
	a.App.Instance = appx.Standardize(m(a.App.Instance))

	a.App.Require(scheduler)
	return a
}

func (a *App) Require(names ...string) *App {
	a.App.Require(names...)
	return a
}
