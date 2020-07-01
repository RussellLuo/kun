package golang

import (
	"github.com/RussellLuo/kok/pkg/appcenter"
)

type Config struct {
	SETTINGS  map[string]interface{}
	APPS      map[string]appcenter.Config
	INSTALLED []string
}

func (c Config) Settings() map[string]interface{} {
	return c.SETTINGS
}

func (c Config) Apps() map[string]appcenter.Config {
	return c.APPS
}

func (c Config) Installed() []string {
	return c.INSTALLED
}
