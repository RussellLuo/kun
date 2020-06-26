package httptest

import (
	"io/ioutil"

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
	StatusCode string `yaml:"statusCode"`
	Body       string `yaml:"body"`
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

type TestSpec struct {
	Imports []string `yaml:"imports"`
	Tests   []test   `yaml:"tests"`
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

	return testSpec, nil
}
