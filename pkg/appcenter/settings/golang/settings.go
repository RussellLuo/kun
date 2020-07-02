package golang

import (
	"github.com/RussellLuo/kok/pkg/appcenter"
)

type Settings struct {
	CONFIG    map[string]interface{}
	APPS      map[string]appcenter.Settings
	INSTALLED []string
}

func (s Settings) Config() appcenter.Config {
	return s.CONFIG
}

func (s Settings) Apps() map[string]appcenter.Settings {
	return s.APPS
}

func (s Settings) Installed() []string {
	return s.INSTALLED
}
