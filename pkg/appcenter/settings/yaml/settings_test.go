package yaml

import (
	"reflect"
	"testing"

	"github.com/RussellLuo/kok/pkg/appcenter"
)

func TestConfig(t *testing.T) {
	cases := []struct {
		inContent     string
		wantConfig    appcenter.Config
		wantApps      map[string]appcenter.Settings
		wantInstalled []string
	}{
		{
			inContent: `---
config:
  ver: v1
apps:
  a:
    config:
      ver: v1
  b:
    config:
      ver: v1
    apps:
      c:
        config:
          ver: v1
    installed:
      - c
installed:
  - a
  - b
`,
			wantConfig: map[string]interface{}{
				"ver": "v1",
			},
			wantApps: map[string]appcenter.Settings{
				"a": Settings{
					CONFIG: map[string]interface{}{
						"ver": "v1",
					},
				},
				"b": Settings{
					CONFIG: map[string]interface{}{
						"ver": "v1",
					},
					APPS: map[string]Settings{
						"c": {
							CONFIG: map[string]interface{}{
								"ver": "v1",
							},
						},
					},
					INSTALLED: []string{"c"},
				},
			},
			wantInstalled: []string{"a", "b"},
		},
	}

	for _, c := range cases {
		settings, err := FromYAML([]byte(c.inContent))
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if !reflect.DeepEqual(settings.Config(), c.wantConfig) {
			t.Fatalf("Config: Got (%#v) != Want (%#v)", settings.Config(), c.wantConfig)
		}

		if !reflect.DeepEqual(settings.Apps(), c.wantApps) {
			t.Fatalf("Apps: Got (%#v) != Want (%#v)", settings.Apps(), c.wantApps)
		}

		if !reflect.DeepEqual(settings.Installed(), c.wantInstalled) {
			t.Fatalf("Installed: Got (%#v) != Want (%#v)", settings.Installed(), c.wantInstalled)
		}
	}
}
