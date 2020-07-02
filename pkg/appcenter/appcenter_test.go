package appcenter

import (
	"context"
	"reflect"
	"testing"
)

type mockSettings struct {
	CONFIG    map[string]interface{}
	APPS      map[string]Settings
	INSTALLED []string
}

func (s mockSettings) Config() Config {
	return s.CONFIG
}

func (s mockSettings) Apps() map[string]Settings {
	return s.APPS
}

func (s mockSettings) Installed() []string {
	return s.INSTALLED
}

func TestInstallRoot(t *testing.T) {
	cases := []struct {
		name             string
		inSettings       Settings
		inRegisteredApps map[string]NewFunc
		wantApp          *App
		wantErr          error
		wantErrStr       string
	}{
		{
			name: "app not registered",
			inSettings: mockSettings{
				CONFIG: map[string]interface{}{
					"ver": "v1",
				},
				APPS: map[string]Settings{
					"a": mockSettings{
						CONFIG: map[string]interface{}{
							"ver": "v1",
						},
					},
					"b": mockSettings{
						CONFIG: map[string]interface{}{
							"ver": "v1",
						},
						APPS: map[string]Settings{
							"c": mockSettings{
								CONFIG: map[string]interface{}{
									"ver": "v1",
								},
							},
						},
						INSTALLED: []string{"c"},
					},
				},
				INSTALLED: []string{"a", "b"},
			},
			inRegisteredApps: map[string]NewFunc{
				"root/a": func(ctx context.Context, config Config) (*App, error) {
					return &App{
						Name: "a",
					}, nil
				},
			},
			wantApp:    nil,
			wantErrStr: `no app registered with name "root/b"`,
		},
		{
			name: "settings not found",
			inSettings: mockSettings{
				CONFIG: map[string]interface{}{
					"ver": "v1",
				},
				APPS: map[string]Settings{
					"b": mockSettings{
						CONFIG: map[string]interface{}{
							"ver": "v1",
						},
						APPS: map[string]Settings{
							"c": mockSettings{
								CONFIG: map[string]interface{}{
									"ver": "v1",
								},
							},
						},
						INSTALLED: []string{"c"},
					},
				},
				INSTALLED: []string{"a", "b"},
			},
			inRegisteredApps: map[string]NewFunc{
				"root/a": func(ctx context.Context, config Config) (*App, error) {
					return &App{
						Name: "a",
					}, nil
				},
				"root/b": func(ctx context.Context, config Config) (*App, error) {
					return &App{
						Name: "b",
					}, nil
				},
			},
			wantApp:    nil,
			wantErrStr: `settings of app "a" is not found`,
		},
		{
			name: "ok",
			inSettings: mockSettings{
				CONFIG: map[string]interface{}{
					"ver": "v1",
				},
				APPS: map[string]Settings{
					"a": mockSettings{
						CONFIG: map[string]interface{}{
							"ver": "v1",
						},
					},
					"b": mockSettings{
						CONFIG: map[string]interface{}{
							"ver": "v1",
						},
						APPS: map[string]Settings{
							"c": mockSettings{
								CONFIG: map[string]interface{}{
									"ver": "v1",
								},
							},
						},
						INSTALLED: []string{"c"},
					},
				},
				INSTALLED: []string{"a", "b"},
			},
			inRegisteredApps: map[string]NewFunc{
				"root/a": func(ctx context.Context, config Config) (*App, error) {
					return &App{
						Name: "a",
					}, nil
				},
				"root/b": func(ctx context.Context, config Config) (*App, error) {
					return &App{
						Name: "b",
					}, nil
				},
				"root/b/c": func(ctx context.Context, config Config) (*App, error) {
					return &App{
						Name: "c",
					}, nil
				},
			},
			wantApp: &App{
				Name: "root",
				subApps: []*App{
					{
						Name: "a",
					},
					{
						Name: "b",
						subApps: []*App{
							{
								Name: "c",
							},
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for name, installFunc := range c.inRegisteredApps {
				if err := Register(name, installFunc); err != nil {
					t.Fatalf("err: %v", err)
				} else {
					defer Unregister(name)
				}
			}

			appName := "root"
			newApp := func(ctx context.Context, config Config) (*App, error) {
				return &App{
					Name: appName,
				}, nil
			}
			app, err := InstallRoot(context.Background(), c.inSettings, appName, newApp)
			if err == nil {
				defer func() {
					if err := app.Uninstall(); err != nil {
						t.Fatalf("err: %v", err)
					}
				}()
			}

			if c.wantErrStr != "" {
				if err == nil || err.Error() != c.wantErrStr {
					t.Fatalf("Err: Got (%#v) != Want (%#v)", err.Error(), c.wantErrStr)
				}
			} else if err != c.wantErr {
				t.Fatalf("Err: Got (%#v) != Want (%#v)", err, c.wantErr)
			}

			if !reflect.DeepEqual(app, c.wantApp) {
				t.Fatalf("App: Got (%#v) != Want (%#v)", app, c.wantApp)
			}
		})
	}
}
