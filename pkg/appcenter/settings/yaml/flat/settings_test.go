package yaml

import (
	"reflect"
	"testing"

	"github.com/RussellLuo/kok/pkg/appcenter/settings/golang"
)

func TestConfig(t *testing.T) {
	cases := []struct {
		inContent     string
		wantConfig    map[string]interface{}
		wantApps      map[string]golang.Settings
		wantInstalled []string
	}{
		{
			inContent: `---
root:
  ver: v1

root/app1:
  ver: v2

root/app2:
  ver: v3

root/app2/app3:
  ver: v4

installed:
  - root/app2
  - root/app2/app3
  - root/app1
  - root
`,
			wantConfig: map[string]interface{}{
				"ver": "v1",
			},
			wantApps: map[string]golang.Settings{
				"app1": {
					CONFIG: map[string]interface{}{
						"ver": "v2",
					},
				},
				"app2": {
					CONFIG: map[string]interface{}{
						"ver": "v3",
					},
					APPS: map[string]golang.Settings{
						"app3": {
							CONFIG: map[string]interface{}{
								"ver": "v4",
							},
						},
					},
					INSTALLED: []string{"app3"},
				},
			},
			wantInstalled: []string{"app2", "app1"},
		},
	}

	for _, c := range cases {
		settings, err := FromYAML([]byte(c.inContent), "installed", "root")
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
