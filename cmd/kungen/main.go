package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RussellLuo/kok/gen"
)

type userFlags struct {
	outDir        string
	flatLayout    bool
	testFileName  string
	formatted     bool
	snakeCase     bool
	enableTracing bool
	oldAnnotation bool

	args []string
}

func main() {
	var flags userFlags
	flag.StringVar(&flags.outDir, "out", ".", "output directory")
	flag.BoolVar(&flags.flatLayout, "flat", true, "whether to use flat layout")
	flag.StringVar(&flags.testFileName, "test", "./http.test.yaml", "the YAML file that provides test-cases for HTTP")
	flag.BoolVar(&flags.formatted, "fmt", true, "whether to make code formatted")
	flag.BoolVar(&flags.snakeCase, "snake", true, "whether to use snake-case for default names")
	flag.BoolVar(&flags.enableTracing, "trace", false, "whether to enable tracing")
	flag.BoolVar(&flags.oldAnnotation, "old", false, "whether to use the old annotation syntax")

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

	generator := gen.New(&gen.Options{
		OutDir:        flags.outDir,
		FlatLayout:    flags.flatLayout,
		SchemaPtr:     true,
		SchemaTag:     "json",
		SnakeCase:     flags.snakeCase,
		Formatted:     flags.formatted,
		EnableTracing: flags.enableTracing,
		OldAnnotation: flags.oldAnnotation,
	})
	files, err := generator.Generate(srcFilename, interfaceName, flags.testFileName)
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := f.Write(); err != nil {
			return err
		}
	}

	return nil
}
