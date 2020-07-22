package yaml

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"

	"github.com/RussellLuo/kok/pkg/appcenter/settings/golang"
	"github.com/RussellLuo/kok/pkg/appcenter/settings/yaml/flat/trie"
)

type Settings map[string]interface{}

func (s Settings) getInstalled(key string) (installed []string, err error) {
	value, ok := s[key]
	if !ok {
		return nil, fmt.Errorf("installed key %q is not found", key)
	}

	err = mapstructure.Decode(value, &installed)
	return
}

func (s Settings) getAppSettings(appName string) (*golang.Settings, error) {
	settings := &golang.Settings{}

	config, ok := s[appName]
	if !ok {
		return settings, nil
	}

	if err := mapstructure.Decode(config, &settings.CONFIG); err != nil {
		return nil, err
	}

	return settings, nil
}

func makeSettings(root string, r *trie.Trie) (settings golang.Settings, err error) {
	r.Walk(func(parent, n *trie.Node) bool {
		if parent == nil || parent.Value == nil {
			return true
		}

		pSettings := parent.Value.(*golang.Settings)
		if pSettings.APPS == nil {
			pSettings.APPS = map[string]golang.Settings{}
		}

		key := strings.TrimPrefix(n.Key, parent.Key+"/")
		pSettings.APPS[key] = *n.Value.(*golang.Settings)
		pSettings.INSTALLED = append(pSettings.INSTALLED, key)

		return false
	})

	value := r.Get(root)
	if value == nil {
		return
	}

	settings = *value.(*golang.Settings)
	return
}

func FromYAML(in []byte, installedKey, rootAppName string) (settings golang.Settings, err error) {
	var s Settings
	if err = yaml.Unmarshal(in, &s); err != nil {
		return
	}

	installed, err := s.getInstalled(installedKey)
	if err != nil {
		return settings, err
	}

	r := trie.New()
	for _, appName := range installed {
		appSettings, err := s.getAppSettings(appName)
		if err != nil {
			return settings, err
		}
		r.Insert(appName, appSettings)
	}

	return makeSettings(rootAppName, r)
}
