package golang

import (
	"github.com/RussellLuo/kok/pkg/appcenter"
)

type Settings struct {
	CONFIG    map[string]interface{}
	APPS      map[string]Settings
	INSTALLED []string
}

func (s Settings) Config() appcenter.Config {
	return s.CONFIG
}

func (s Settings) Apps() map[string]appcenter.Settings {
	apps := make(map[string]appcenter.Settings)
	for name, settings := range s.APPS {
		apps[name] = settings
	}
	return apps
}

func (s Settings) Installed() []string {
	return s.INSTALLED
}
