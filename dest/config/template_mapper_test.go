package config

import (
	"os"
	"testing"

	"github.com/tlmiller/disttrust/dest"
)

func TestTemplateMapperFailsForNoData(t *testing.T) {
	if _, err := TemplateMapper(NewTemplate()); err == nil {
		t.Fatal("expected error for template mapper with no config data")
	}
}

func TestTemplateMapperOutputMeta(t *testing.T) {
	conf := &Template{
		Gid:    "group",
		Source: "source-file.tmpl",
		Mode:   "0700",
		Out:    "out-file.txt",
		Uid:    "user",
	}

	d, err := TemplateMapper(conf)
	if err != nil {
		t.Fatalf("unexpected error getting template dest from mapper: %v", err)
	}

	tf, ok := d.(*dest.TemplateFile)
	if !ok {
		t.Fatal("unexpected type from template file mapper")
	}

	if tf.Dest.Dest.Gid != "group" {
		t.Errorf("template file dest group %s != group", tf.Dest.Dest.Gid)
	}
	if tf.Dest.Dest.Path != "out-file.txt" {
		t.Errorf("template file dest path %s != out-file.txt", tf.Dest.Dest.Path)
	}
	if tf.Dest.Dest.Mode != os.FileMode(0700) {
		t.Errorf("template file dest mode %d != 0700", tf.Dest.Dest.Mode)
	}
	if tf.Dest.Dest.Uid != "user" {
		t.Errorf("template file dest user %s != user", tf.Dest.Dest.Uid)
	}
}
