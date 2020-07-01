package appcenter_test

// We use `appcenter_test` as the package name, since the `golang` package is
// imported and used in the tests below, which will cause an "import cycle"
// if the same tests are located in the `appcenter` package.

import (
	"context"
	"reflect"
	"testing"

	"github.com/RussellLuo/kok/pkg/appcenter"
	"github.com/RussellLuo/kok/pkg/appcenter/config/golang"
)

func TestInstall(t *testing.T) {
	cases := []struct {
		name             string
		inConfig         appcenter.Config
		inRegisteredApps map[string]appcenter.InstallFunc
		wantApps         []*appcenter.App
		wantErr          error
		wantErrStr       string
	}{
		{
			name: "app not registered",
			inConfig: golang.Config{
				SETTINGS: map[string]interface{}{
					"ver": "v1",
				},
				APPS: map[string]appcenter.Config{
					"a": golang.Config{
						SETTINGS: map[string]interface{}{
							"ver": "v1",
						},
					},
					"b": golang.Config{
						SETTINGS: map[string]interface{}{
							"ver": "v1",
						},
						APPS: map[string]appcenter.Config{
							"c": golang.Config{
								SETTINGS: map[string]interface{}{
									"ver": "v1",
								},
							},
						},
						INSTALLED: []string{"c"},
					},
				},
				INSTALLED: []string{"a", "b"},
			},
			inRegisteredApps: map[string]appcenter.InstallFunc{
				"a": func(ctx context.Context, name string, config appcenter.Config) (*appcenter.App, error) {
					return nil, nil
				},
			},
			wantApps:   nil,
			wantErrStr: `app "b" is not registered`,
		},
		{
			name: "configuration not found",
			inConfig: golang.Config{
				SETTINGS: map[string]interface{}{
					"ver": "v1",
				},
				APPS: map[string]appcenter.Config{
					"b": golang.Config{
						SETTINGS: map[string]interface{}{
							"ver": "v1",
						},
						APPS: map[string]appcenter.Config{
							"c": golang.Config{
								SETTINGS: map[string]interface{}{
									"ver": "v1",
								},
							},
						},
						INSTALLED: []string{"c"},
					},
				},
				INSTALLED: []string{"a", "b"},
			},
			inRegisteredApps: map[string]appcenter.InstallFunc{
				"a": func(ctx context.Context, name string, config appcenter.Config) (*appcenter.App, error) {
					return nil, nil
				},
				"b": func(ctx context.Context, name string, config appcenter.Config) (*appcenter.App, error) {
					return nil, nil
				},
			},
			wantApps:   nil,
			wantErrStr: `configuration of app "a" is not found`,
		},
		{
			name: "ok",
			inConfig: golang.Config{
				SETTINGS: map[string]interface{}{
					"ver": "v1",
				},
				APPS: map[string]appcenter.Config{
					"a": golang.Config{
						SETTINGS: map[string]interface{}{
							"ver": "v1",
						},
					},
					"b": golang.Config{
						SETTINGS: map[string]interface{}{
							"ver": "v1",
						},
						APPS: map[string]appcenter.Config{
							"c": golang.Config{
								SETTINGS: map[string]interface{}{
									"ver": "v1",
								},
							},
						},
						INSTALLED: []string{"c"},
					},
				},
				INSTALLED: []string{"a", "b"},
			},
			inRegisteredApps: map[string]appcenter.InstallFunc{
				"a": func(ctx context.Context, name string, config appcenter.Config) (*appcenter.App, error) {
					return &appcenter.App{
						Name: "a",
					}, nil
				},
				"b": func(ctx context.Context, name string, config appcenter.Config) (*appcenter.App, error) {
					return &appcenter.App{
						Name: "b",
					}, nil
				},
			},
			wantApps: []*appcenter.App{
				{
					Name: "a",
				},
				{
					Name: "b",
				},
			},
			wantErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for name, installFunc := range c.inRegisteredApps {
				appcenter.Register(name, installFunc)
				defer appcenter.Unregister(name)
			}

			apps, err := appcenter.Install(context.Background(), c.inConfig.Installed(), c.inConfig.Apps())

			if c.wantErrStr != "" {
				if err == nil || err.Error() != c.wantErrStr {
					t.Fatalf("Err: Got (%#v) != Want (%#v)", err.Error(), c.wantErrStr)
				}
			} else if err != c.wantErr {
				t.Fatalf("Err: Got (%#v) != Want (%#v)", err, c.wantErr)
			}

			if !reflect.DeepEqual(apps, c.wantApps) {
				t.Fatalf("Apps: Got (%#v) != Want (%#v)", apps, c.wantApps)
			}
		})
	}
}
