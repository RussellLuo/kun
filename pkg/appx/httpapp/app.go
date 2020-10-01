package httpapp

import (
	"github.com/RussellLuo/appx"
)

type App struct {
	*appx.App

	parent  string
	pattern string
}

func New(name string) *App {
	return &App{App: appx.New(name)}
}

func (a *App) MountOn(parent string) *App {
	a.parent = parent
	a.App.Require(parent)
	return a
}

func (a *App) Pattern(pattern string) *App {
	a.pattern = pattern
	return a
}

func (a *App) Require(names ...string) *App {
	a.App.Require(names...)
	return a
}

func (a *App) InitFunc(initFunc appx.InitFuncV2) *appx.App {
	init := initFunc
	if a.parent != "" {
		init = a.mountOnParent(initFunc)
	}

	a.App.InitFunc(init)
	return a.App // Return the wrapped *appx.App
}

func (a *App) mountOnParent(initFunc appx.InitFuncV2) appx.InitFuncV2 {
	return func(ctx appx.Context) error {
		if err := initFunc(ctx); err != nil {
			return err
		}

		parent, err := GetRouter(ctx.Required[a.parent].Value)
		if err != nil {
			return err
		}

		r, err := GetRouter(ctx.App.Value)
		if err != nil {
			return err
		}

		MountRouter(parent, a.pattern, r)
		return nil
	}
}
