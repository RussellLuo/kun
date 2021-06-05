package wrapper

import (
	"context"

	"github.com/RussellLuo/appx"
)

type InitWrapper struct {
	instance appx.Initializer

	afterInitFunc appx.InitFuncV2
}

func New(instance appx.Initializer) *InitWrapper {
	return &InitWrapper{instance: instance}
}

func (w *InitWrapper) AfterInitFunc(f appx.InitFuncV2) {
	w.afterInitFunc = f
}

func (w *InitWrapper) Init(ctx appx.Context) error {
	if err := w.instance.Init(ctx); err != nil {
		return err
	}
	if w.afterInitFunc != nil {
		if err := w.afterInitFunc(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (w *InitWrapper) Clean() error {
	if cleaner, ok := w.instance.(appx.Cleaner); ok {
		return cleaner.Clean()
	}
	return nil
}

func (w *InitWrapper) Start(ctx context.Context) error {
	if startStopper, ok := w.instance.(appx.StartStopper); ok {
		return startStopper.Start(ctx)
	}
	return nil
}

func (w *InitWrapper) Stop(ctx context.Context) error {
	if startStopper, ok := w.instance.(appx.StartStopper); ok {
		return startStopper.Stop(ctx)
	}
	return nil
}

func (w *InitWrapper) Instance() interface{} {
	return w.instance
}

func (w *InitWrapper) Validate() error {
	if validator, ok := w.instance.(appx.Validator); ok {
		return validator.Validate()
	}
	return nil
}
