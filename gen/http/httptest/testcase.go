package httptest

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type request struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
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

func getTests(testFilename string) (tests []test, err error) {
	b, err := ioutil.ReadFile(testFilename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(b), &tests)
	if err != nil {
		return nil, err
	}

	return tests, nil
}
