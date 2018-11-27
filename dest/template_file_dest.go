package dest

import (
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/file"
	"github.com/tlmiller/disttrust/provider"
)

type fileWrapper struct {
	Dest     file.File
	origPath string
	ofile    *os.File
}

type TemplateFile struct {
	Loader TemplateLoader
	Dest   *fileWrapper
}

func (t *fileWrapper) Close() error {
	if t.ofile == nil {
		return nil
	}
	err := t.ofile.Close()
	if err != nil {
		t.ofile = nil
		return err
	}
	err = t.Dest.Chown()
	t.ofile = nil
	return err
}

func NewTemplateFile(loader TemplateLoader, dest file.File) *TemplateFile {
	return &TemplateFile{
		Loader: loader,
		Dest: &fileWrapper{
			Dest:     dest,
			origPath: dest.Path,
		},
	}
}

func (t *TemplateFile) Send(res *provider.Response) error {
	defer t.Dest.Close()
	tmplPath := strings.Builder{}
	if err := NewTemplate(TemplateString(t.Dest.origPath), &tmplPath).Send(res); err != nil {
		return errors.Wrapf(err, "parsing possible template dest path %s", t.Dest.origPath)
	}
	t.Dest.Dest.Path = tmplPath.String()
	return NewTemplate(t.Loader, t.Dest).Send(res)
}

func (t *fileWrapper) Write(p []byte) (n int, err error) {
	if t.ofile == nil {
		of, err := os.OpenFile(t.Dest.Path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			t.Dest.Mode)
		if err != nil {
			return 0, err
		}
		t.ofile = of
	}
	return t.ofile.Write(p)
}
