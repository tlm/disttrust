package dest

import (
	"io"
	"io/ioutil"
	"text/template"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/provider"
)

type TemplateLoader interface {
	Load() (string, error)
}

type TemplateLoaderFunc func() (string, error)

func TemplateFileLoader(path string) TemplateLoader {
	return TemplateLoaderFunc(func() (string, error) {
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return "", errors.Wrap(err, "reading template file")
		}
		return string(buf), nil
	})
}

type TemplateString string

type Template struct {
	Loader TemplateLoader
	Dest   io.Writer
}

func (f TemplateLoaderFunc) Load() (string, error) {
	return f()
}

func (s TemplateString) Load() (string, error) {
	return string(s), nil
}

func NewTemplate(loader TemplateLoader, dest io.Writer) *Template {
	return &Template{
		Loader: loader,
		Dest:   dest,
	}
}

func (t *Template) Send(res *provider.Response) error {
	tmplBody, err := t.Loader.Load()
	if err != nil {
		return errors.Wrap(err, "loading template")
	}

	tmpl, err := template.New("template").Parse(tmplBody)
	if err != nil {
		return errors.Wrap(err, "parsing dest template")
	}

	err = tmpl.Execute(t.Dest, res)
	if err != nil {
		return errors.Wrap(err, "writing dest template")
	}
	return nil
}
