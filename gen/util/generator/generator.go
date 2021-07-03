package generator

import (
	"bytes"
	"text/template"
)

var formatters = []Formatter{Gofmt, Goimports}

type File struct {
	Name    string
	Content []byte
}

type PkgInfo struct {
	CurrentPkgName    string
	EndpointPkgPrefix string
	EndpointPkgPath   string
}

type Options struct {
	Name           string
	Funcs          template.FuncMap
	Formatted      bool
	TargetFileName string
}

func Generate(text string, data interface{}, opts Options) (*File, error) {
	tmpl, err := template.New(opts.Name).Funcs(opts.Funcs).Parse(text)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	b := buf.Bytes()

	if opts.Formatted {
		for _, fmt := range formatters {
			if b, err = fmt(b); err != nil {
				return nil, err
			}
		}
	}

	return &File{
		Name:    opts.TargetFileName,
		Content: b,
	}, nil
}
