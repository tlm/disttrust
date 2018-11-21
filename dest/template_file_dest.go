package dest

import (
	"os"

	"github.com/tlmiller/disttrust/file"
	"github.com/tlmiller/disttrust/provider"
)

type fileWrapper struct {
	Dest  file.File
	ofile *os.File
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
		return err
	}
	err = t.Dest.Chown()
	return err
}

func NewTemplateFile(loader TemplateLoader, dest file.File) *TemplateFile {
	return &TemplateFile{
		Loader: loader,
		Dest: &fileWrapper{
			Dest: dest,
		},
	}
}

func (t *TemplateFile) Send(res *provider.Response) error {
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
