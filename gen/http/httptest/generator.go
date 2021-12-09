package httptest

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/RussellLuo/kok/gen/util/annotation"
	"github.com/RussellLuo/kok/gen/util/generator"
	"github.com/RussellLuo/kok/pkg/ifacetool"
)

var (
	template = annotation.FileHeader + `
package {{.PkgInfo.CurrentPkgName}}

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	{{- range .Data.Imports}}
	{{.ImportString}}
	{{- end}}

	{{- range .TestSpec.Imports}}
	{{.Alias}} {{.Path}}
	{{- end}}
)

{{- $srcPkgPrefix := .Data.SrcPkgQualifier}}
{{- $interfaceName := .Data.InterfaceName}}
{{- $mockInterfaceName := printf "%s%s" $interfaceName "Mock"}}

// Ensure that {{$mockInterfaceName}} does implement {{$srcPkgPrefix}}{{$interfaceName}}.
var _ {{$srcPkgPrefix}}{{$interfaceName}} = &{{$mockInterfaceName}}{}

type {{$mockInterfaceName}} struct {
{{- range .Data.Methods}}
	{{.Name}}Func func({{.ArgList}}) {{.ReturnArgNamedValueList}}
{{- end}}
}

{{- range .Data.Methods}}

func (mock *{{$mockInterfaceName}}) {{.Name}}({{.ArgList}}) {{.ReturnArgNamedValueList}} {
	if mock.{{.Name}}Func == nil {
		panic("{{$mockInterfaceName}}.{{.Name}}Func: not implemented")
	}
	return mock.{{.Name}}Func({{.CallArgList}})
}
{{- end}}

type request struct {
	method string
	path   string
	header map[string]string
	body   string
}

func (r request) ServedBy(handler http.Handler) *httptest.ResponseRecorder {
	var req *http.Request
	if r.body != "" {
		reqBody := strings.NewReader(r.body)
		req = httptest.NewRequest(r.method, r.path, reqBody)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	} else {
		req = httptest.NewRequest(r.method, r.path, nil)
	}

	for key, value := range r.header {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	return w
}

type response struct {
	statusCode  int
	contentType string
	body        []byte
}

func (want response) Equal(w *httptest.ResponseRecorder) string {
	resp := w.Result()
	gotBody, _ := ioutil.ReadAll(resp.Body)

	gotStatusCode := resp.StatusCode
	if gotStatusCode != want.statusCode {
		return fmt.Sprintf("StatusCode: got (%d), want (%d)", gotStatusCode, want.statusCode)
	}

	wantContentType := want.contentType
	if wantContentType == "" {
		wantContentType = "application/json; charset=utf-8"
	}

	gotContentType := resp.Header.Get("Content-Type")
	if gotContentType != wantContentType {
		return fmt.Sprintf("ContentType: got (%q), want (%q)", gotContentType, wantContentType)
	}

	if strings.HasPrefix(gotContentType, "application/json") {
		// Remove the trailing newline from the JSON bytes encoded by Go.
		// See https://github.com/golang/go/issues/37083.
		gotBody = bytes.TrimSuffix(gotBody, []byte("\n"))
	}

	if !bytes.Equal(gotBody, want.body) {
		return fmt.Sprintf("Body: got (%q), want (%q)", gotBody, want.body)
	}

	return ""
}

{{- $codecs := .TestSpec.Codecs}}
{{- range .TestSpec.Tests}}

{{$method := method .Name}}
{{$params := $method.Params}}
{{$returns := $method.Returns}}
{{$nonCtxParams := nonCtxParams $params}}

func TestHTTP_{{.Name}}(t *testing.T) {
	// in contains all the input parameters (except ctx) of {{.Name}}.
	type in struct {
		{{- range $nonCtxParams}}
		{{.Name}} {{.TypeString}}
		{{- end}}
	}

	// out contains all the output parameters of {{.Name}}.
	type out struct {
		{{- range $returns}}
		{{.Name}} {{.TypeString}}
		{{- end}}
	}

	{{if .Cases -}}
	cases := []struct {
		name         string
		request      request
		wantIn       in
		out          out
		wantResponse response
	}{
		{{- range .Cases}}
		{
			name: "{{.Name}}",
			request: request{
				method: "{{.Request.Method}}",
				path:   "{{.Request.Path}}",
				{{- if .Request.Header}}
				header: map[string]string{
					{{- range $key, $value := .Request.Header}}
					"{{$key}}": "{{$value}}",
					{{- end}}
				},
				{{- end}}
				{{- if .Request.Body}}
				body:   ` + "`{{.Request.Body}}`" + `,
				{{- end}}
			},
			wantIn: in{
				{{.WantIn}}
			},
			out: out{
				{{.Out}}
			},
			wantResponse: response{
				statusCode: {{.WantResponse.StatusCode}},
				{{- if .WantResponse.ContentType}}
				contentType: "{{.WantResponse.ContentType}}",
				{{- end}}
				{{- if .WantResponse.Body}}
				body: {{bodyToBytes .WantResponse.Body}},
				{{- end}}
			},
		},
		{{- end}}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotIn in
			w := c.request.ServedBy(NewHTTPRouter(
				&{{$mockInterfaceName}}{
					{{.Name}}Func: func({{$method.ArgList}}) {{$method.ReturnArgNamedValueList}} {
						gotIn = in{
							{{- range $nonCtxParams}}
							{{.Name}}: {{.Name}},
							{{- end}}
						}
						return {{joinParams $returns "c.out.$Name" ", "}}
					},
				},
				{{$codecs}},
			))

			if !reflect.DeepEqual(gotIn, c.wantIn) {
				t.Fatalf("In: got (%v), want (%v)", gotIn, c.wantIn)
			}

			if errStr := c.wantResponse.Equal(w); errStr != "" {
				t.Fatal(errStr)
			}
		})
	}
	{{- end}}
}
{{- end}}
`
)

type Options struct {
	Formatted bool
}

type Generator struct {
	opts *Options
}

func New(opts *Options) *Generator {
	return &Generator{opts: opts}
}

func (g *Generator) Generate(pkgInfo *generator.PkgInfo, ifaceData *ifacetool.Data, testFilename string) (*generator.File, error) {
	testSpec, err := getTestSpec(testFilename)
	if err != nil {
		return nil, err
	}

	data := struct {
		PkgInfo  *generator.PkgInfo
		Data     *ifacetool.Data
		TestSpec *TestSpec
	}{
		PkgInfo:  pkgInfo,
		Data:     ifaceData,
		TestSpec: testSpec,
	}

	methodMap := make(map[string]*ifacetool.Method)
	for _, method := range ifaceData.Methods {
		methodMap[method.Name] = method
	}

	return generator.Generate(template, data, generator.Options{
		Funcs: map[string]interface{}{
			"joinParams": func(params []*ifacetool.Param, format, sep string) string {
				var results []string

				for _, p := range params {
					r := strings.NewReplacer("$Name", p.Name)
					results = append(results, r.Replace(format))
				}
				return strings.Join(results, sep)
			},
			"method": func(name string) *ifacetool.Method {
				method, ok := methodMap[name]
				if !ok {
					return nil
				}
				return method
			},
			"methodParams": func(name string) []*ifacetool.Param {
				method, ok := methodMap[name]
				if !ok {
					return nil
				}
				return method.Params
			},
			"methodReturns": func(name string) []*ifacetool.Param {
				method, ok := methodMap[name]
				if !ok {
					return nil
				}
				return method.Returns
			},
			"nonCtxParams": func(params []*ifacetool.Param) (out []*ifacetool.Param) {
				for _, p := range params {
					if p.TypeString != "context.Context" {
						out = append(out, p)
					}
				}
				return
			},
			"bodyToBytes": func(s string) string {
				if s == "" {
					// An empty string indicates a nil byte slice.
					return "[]byte(nil)"
				}

				if strings.HasPrefix(s, "0x") {
					// This is a hexadecimal string, decode it into bytes.
					//
					// Note that kok borrows the idea from eth2.0 to represent binary data
					// as hex encoded strings, see https://github.com/ethereum/eth2.0-spec-tests/issues/5.
					decoded, err := hex.DecodeString(s[2:])
					if err != nil {
						panic(err)
					}

					var hexes []string
					for _, b := range decoded {
						hexes = append(hexes, fmt.Sprintf("0x%x", b))
					}
					return fmt.Sprintf("[]byte{%s}", strings.Join(hexes, ", "))
				}

				// This is a normal string, leave it as is.
				return fmt.Sprintf("[]byte(`%s`)", s)
			},
		},
		Formatted:      g.opts.Formatted,
		TargetFileName: "http_test.go",
	})
}
