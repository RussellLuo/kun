package generator

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"text/template"
)

var formatters = []Formatter{Gofmt, Goimports}

type File struct {
	Name    string
	Content []byte
}

func (f *File) MoveTo(dir string) {
	f.Name = filepath.Join(dir, f.Name)
}

func (f *File) Write() error {
	return ioutil.WriteFile(f.Name, f.Content, 0644)
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
