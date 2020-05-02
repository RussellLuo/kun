package gen

import (
	"bytes"
	"text/template"
)

type Options struct {
	Name       string
	Funcs      template.FuncMap
	Formatters []Formatter
}

func Generate(text string, data interface{}, opts Options) ([]byte, error) {
	tmpl, err := template.New(opts.Name).Funcs(opts.Funcs).Parse(text)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	b := buf.Bytes()
	for _, fmt := range opts.Formatters {
		if b, err = fmt(b); err != nil {
			return nil, err
		}
	}

	return b, nil
}
