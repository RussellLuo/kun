package yaml

import (
	"reflect"
	"testing"

	"github.com/RussellLuo/kok/pkg/appcenter"
)

func TestConfig(t *testing.T) {
	cases := []struct {
		inContent     string
		wantSettings  map[string]interface{}
		wantApps      map[string]appcenter.Config
		wantInstalled []string
	}{
		{
			inContent: `---
settings:
  ver: v1
apps:
  a:
    settings:
      ver: v1
  b:
    settings:
      ver: v1
    apps:
      c:
        settings:
          ver: v1
    installed:
      - c
installed:
  - a
  - b
`,
			wantSettings: map[string]interface{}{
				"ver": "v1",
			},
			wantApps: map[string]appcenter.Config{
				"a": Config{
					SETTINGS: map[string]interface{}{
						"ver": "v1",
					},
				},
				"b": Config{
					SETTINGS: map[string]interface{}{
						"ver": "v1",
					},
					APPS: map[string]Config{
						"c": {
							SETTINGS: map[string]interface{}{
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
		config, err := FromYAML([]byte(c.inContent))
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if !reflect.DeepEqual(config.Settings(), c.wantSettings) {
			t.Fatalf("Settings: Got (%#v) != Want (%#v)", config.Settings(), c.wantSettings)
		}

		if !reflect.DeepEqual(config.Apps(), c.wantApps) {
			t.Fatalf("Apps: Got (%#v) != Want (%#v)", config.Apps(), c.wantApps)
		}

		if !reflect.DeepEqual(config.Installed(), c.wantInstalled) {
			t.Fatalf("Installed: Got (%#v) != Want (%#v)", config.Installed(), c.wantInstalled)
		}
	}
}
