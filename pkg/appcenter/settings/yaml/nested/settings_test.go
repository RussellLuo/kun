package yaml

import (
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	cases := []struct {
		inContent     string
		wantConfig    map[string]interface{}
		wantApps      map[string]Settings
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
			wantApps: map[string]Settings{
				"a": {
					CONFIG: map[string]interface{}{
						"ver": "v1",
					},
				},
				"b": {
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

		if !reflect.DeepEqual(settings.CONFIG, c.wantConfig) {
			t.Fatalf("Config: Got (%#v) != Want (%#v)", settings.CONFIG, c.wantConfig)
		}

		if !reflect.DeepEqual(settings.APPS, c.wantApps) {
			t.Fatalf("Apps: Got (%#v) != Want (%#v)", settings.APPS, c.wantApps)
		}

		if !reflect.DeepEqual(settings.INSTALLED, c.wantInstalled) {
			t.Fatalf("Installed: Got (%#v) != Want (%#v)", settings.INSTALLED, c.wantInstalled)
		}
	}
}
