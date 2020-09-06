package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/RussellLuo/kok/gen"
)

type userFlags struct {
	outDir        string
	pkgName       string
	testFileName  string
	formatted     bool
	enableTracing bool

	args []string
}

func main() {
	var flags userFlags
	flag.StringVar(&flags.outDir, "out", ".", "output directory")
	flag.StringVar(&flags.pkgName, "pkg", "", "package name (default will infer)")
	flag.StringVar(&flags.testFileName, "test", "./http.test.yaml", "the YAML file that provides test-cases for HTTP")
	flag.BoolVar(&flags.formatted, "fmt", true, "whether to make code formatted")
	flag.BoolVar(&flags.enableTracing, "trace", false, "whether to enable tracing")

	flag.Usage = func() {
		fmt.Println(`kokgen [flags] source-file interface-name`)
		flag.PrintDefaults()
	}

	flag.Parse()
	flags.args = flag.Args()

	if err := run(flags); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}
}

func run(flags userFlags) error {
	if len(flags.args) != 2 {
		return errors.New("need 2 arguments")
	}

	srcFilename, interfaceName := flags.args[0], flags.args[1]

	srcFilename, err := filepath.Abs(srcFilename)
	if err != nil {
		return err
	}

	content, err := gen.New(gen.Options{
		SchemaPtr:         true,
		SchemaTag:         "json",
		TagKeyToSnakeCase: true,
		Formatted:         flags.formatted,
		EnableTracing:     flags.enableTracing,
	}).Generate(srcFilename, interfaceName, flags.pkgName, flags.testFileName)
	if err != nil {
		return err
	}

	if flags.outDir != "." {
		if err = os.MkdirAll(flags.outDir, 0755); err != nil {
			return err
		}
	}

	files := map[string][]byte{
		"endpoint.go":    content.Endpoint,
		"http.go":        content.HTTP,
		"http_test.go":   content.HTTPTest,
		"http_client.go": content.HTTPClient,
		"oasv2.go":       content.OASv2,
	}
	for name, data := range files {
		if len(data) == 0 {
			continue
		}

		if err := ioutil.WriteFile(filepath.Join(flags.outDir, name), data, 0644); err != nil {
			return err
		}
	}

	return nil
}
