package appcenter

import (
	"context"
	"fmt"
)

var (
	registry = make(map[string]InstallFunc)
)

// Config the configuration of an application.
type Config interface {
	// Settings returns the corresponding settings of the current application.
	Settings() map[string]interface{}

	// Apps returns the configurations of all sub-applications.
	Apps() map[string]Config

	// Installed returns a list of sub-application names that will be installed
	// within the current application.
	Installed() []string
}

// InstallFunc installs an application. It will return an error if failed.
type InstallFunc func(ctx context.Context, name string, config Config) (*App, error)

// UninstallFunc uninstalls the corresponding application. It will return an error if failed.
type UninstallFunc func() error

// App represents an application.
type App struct {
	Name      string
	Uninstall UninstallFunc

	// The use-case specific options, which are customized by users.
	Options interface{}
}

// Register registers an application with its name and the install function.
func Register(name string, install InstallFunc) error {
	if _, ok := registry[name]; ok {
		return fmt.Errorf("app %q already exists", name)
	}
	registry[name] = install
	return nil
}

// Unregister unregisters the given applications specified by names. It will
// clear the registry (i.e. unregister all applications) if no name is provided.
func Unregister(names ...string) {
	if len(names) == 0 {
		// Clear the registry.
		registry = make(map[string]InstallFunc)
	}

	for _, name := range names {
		delete(registry, name)
	}
}

func Install(ctx context.Context, names []string, configs map[string]Config) (apps []*App, err error) {
	for _, name := range names {
		install, ok := registry[name]
		if !ok {
			return nil, fmt.Errorf("app %q is not registered", name)
		}

		config, ok := configs[name]
		if !ok {
			return nil, fmt.Errorf("configuration of app %q is not found", name)
		}

		app, err := install(ctx, name, config)
		if err != nil {
			return nil, err
		}

		apps = append(apps, app)
	}
	return
}
