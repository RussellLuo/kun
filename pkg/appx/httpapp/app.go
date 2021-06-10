package httpapp

import (
	"github.com/RussellLuo/appx"
	"github.com/RussellLuo/kok/pkg/appx/wrapper"
)

type App struct {
	*appx.App

	wrapper *wrapper.InitWrapper

	parent  string
	pattern string
}

func New(name string, instance appx.Initializer) *App {
	w := wrapper.New(instance)
	return &App{
		App:     appx.New(name, w),
		wrapper: w,
	}
}

func (a *App) MountOn(parent string, pattern string) *App {
	a.parent = parent
	a.App.Require(parent)

	a.pattern = pattern
	a.wrapper.AfterInitFunc(a.afterInit)

	return a
}

func (a *App) Require(names ...string) *App {
	a.App.Require(names...)
	return a
}

func (a *App) afterInit(ctx appx.Context) error {
	parent, err := GetChiRouter(ctx.MustLoad(a.parent))
	if err != nil {
		return err
	}

	r, err := GetChiRouter(a.wrapper.Instance())
	if err != nil {
		return err
	}

	MountRouter(parent, a.pattern, r)
	return nil
}
