package httptest

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type header map[string]string

type request struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
	Header header `yaml:"header"`
	Body   string `yaml:"body"`
}

type response struct {
	StatusCode  string `yaml:"statusCode"`
	ContentType string `yaml:"contentType"`
	Body        string `yaml:"body"`
}

type testCase struct {
	Name         string   `yaml:"name"`
	Request      request  `yaml:"request"`
	WantIn       string   `yaml:"wantIn"`
	Out          string   `yaml:"out"`
	WantResponse response `yaml:"wantResponse"`
}

type test struct {
	Name  string     `yaml:"name"`
	Cases []testCase `yaml:"cases"`
}

type Import struct {
	Path  string `yaml:"path"`
	Alias string `yaml:"alias"`
}

type TestSpec struct {
	RawImports []string `yaml:"imports"`
	Imports    []Import `yaml:"-"`
	Codecs     string   `yaml:"codecs"`
	Tests      []test   `yaml:"tests"`
}

func getTestSpec(testFilename string) (*TestSpec, error) {
	b, err := ioutil.ReadFile(testFilename)
	if err != nil {
		return nil, err
	}

	testSpec := &TestSpec{}
	err = yaml.Unmarshal([]byte(b), testSpec)
	if err != nil {
		return nil, err
	}

	imports := getImports(testSpec.RawImports)
	testSpec.Imports = append(testSpec.Imports, imports...)

	return testSpec, nil
}

func getImports(rawImports []string) (imports []Import) {
	var path, alias string

	for i, str := range rawImports {
		fields := strings.Fields(str)
		switch len(fields) {
		case 1:
			alias, path = "", fields[0]
		case 2:
			alias, path = fields[0], fields[1]
		default:
			panic(fmt.Errorf("invalid path in imports[%d]: %s", i, str))
		}

		if !strings.HasPrefix(path, `"`) {
			path = fmt.Sprintf("%q", path)
		}
		imports = append(imports, Import{
			Path:  path,
			Alias: alias,
		})
	}

	return
}
