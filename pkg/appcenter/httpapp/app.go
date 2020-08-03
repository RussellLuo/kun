package httpapp

import (
	"context"

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

func (a *App) Init(initFunc appx.InitFunc) *appx.App {
	init := initFunc
	if a.parent != "" {
		init = a.mountOnParent(initFunc)
	}

	a.App.Init2(init)
	return a.App // Return the wrapped *appx.App
}

func (a *App) mountOnParent(initFunc appx.InitFunc) appx.InitFunc {
	return func(ctx context.Context, lc appx.Lifecycle, apps map[string]*appx.App) (appx.Value, appx.CleanFunc, error) {
		value, clean, err := initFunc(ctx, lc, apps)
		if err != nil {
			return nil, nil, err
		}

		parent, err := GetRouter(apps[a.parent].Value)
		if err != nil {
			return nil, nil, err
		}

		r, err := GetRouter(value)
		if err != nil {
			return nil, nil, err
		}

		MountRouter(parent, a.pattern, r)
		return value, clean, nil
	}
}
