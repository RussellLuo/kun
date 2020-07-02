package yaml

import (
	"gopkg.in/yaml.v2"

	"github.com/RussellLuo/kok/pkg/appcenter"
)

type Settings struct {
	CONFIG    map[string]interface{} `yaml:"config"`
	APPS      map[string]Settings    `yaml:"apps"`
	INSTALLED []string               `yaml:"installed"`
}

func FromYAML(in []byte) (settings Settings, err error) {
	err = yaml.Unmarshal(in, &settings)
	return
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
