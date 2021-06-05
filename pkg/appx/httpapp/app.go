package httpapp

import (
	"github.com/RussellLuo/appx"
	"github.com/RussellLuo/kok/pkg/appx/wrapper"
	"github.com/go-chi/chi"
)

type App struct {
	*appx.App

	wrapper *wrapper.InitWrapper

	parent  string
	pattern string
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

func (a *App) MountOn(parent string, pattern ...string) *App {
	a.parent = parent
	a.App.Require(parent)

	if len(pattern) == 1 {
		a.pattern = pattern[0]
		a.wrapper.AfterInitFunc(a.afterInit)
	}

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

func (a *App) afterInit(ctx appx.Context) error {
	var parent chi.Router
	var err error

	instance := ctx.MustLoad(a.parent)
	if inst, ok := instance.(ChiRouter); ok {
		parent, err = GetChiRouter(inst)
		if err != nil {
			return err
		}
	} else {
		// The parent app is an old-style one.
		// TODO: remove this snippet when all apps are new-style ones.
		parent, err = GetRouter(ctx.Required[a.parent].Value)
		if err != nil {
			return err
		}
	}

	r, err := GetChiRouter(a.wrapper.Instance())
	if err != nil {
		return err
	}

	MountRouter(parent, a.pattern, r)
	return nil
}
