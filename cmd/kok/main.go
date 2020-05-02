package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/RussellLuo/kok/kok"
)

type userFlags struct {
	outDir  string
	pkgName string
	args    []string
}

func main() {
	var flags userFlags
	flag.StringVar(&flags.outDir, "out-dir", ".", "output directory")
	flag.StringVar(&flags.pkgName, "pkg", "", "package name (default will infer)")

	flag.Usage = func() {
		fmt.Println(`kok [flags] source-file interface-name`)
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

	content, err := kok.New(kok.Options{
		SchemaPtr:         true,
		SchemaTag:         "json",
		TagKeyToSnakeCase: true,
	}).Generate(srcFilename, interfaceName, flags.pkgName)
	if err != nil {
		return err
	}

	if flags.outDir != "." {
		if err = os.MkdirAll(flags.outDir, 0755); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(filepath.Join(flags.outDir, "endpoint.go"), content.Endpoint, 0644); err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(flags.outDir, "http.go"), content.HTTP, 0644)
}
