package yaml

import (
	"gopkg.in/yaml.v2"

	"github.com/RussellLuo/kok/pkg/appcenter"
)

type Config struct {
	SETTINGS  map[string]interface{} `yaml:"settings"`
	APPS      map[string]Config      `yaml:"apps"`
	INSTALLED []string               `yaml:"installed"`
}

func FromYAML(in []byte) (config Config, err error) {
	err = yaml.Unmarshal(in, &config)
	return
}

func (c Config) Settings() map[string]interface{} {
	return c.SETTINGS
}

func (c Config) Apps() map[string]appcenter.Config {
	apps := make(map[string]appcenter.Config)
	for name, config := range c.APPS {
		apps[name] = config
	}
	return apps
}

func (c Config) Installed() []string {
	return c.INSTALLED
}
