package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/tlmiller/disttrust/util"
)

type Builder struct {
	configs []string
	V       *viper.Viper
}

func (b *Builder) BuildAndValidate() error {
	SetDefaults(b.V)
	return b.readConfigs()
}

func NewBuilder(configs []string) *Builder {
	return &Builder{
		configs: configs,
		V:       viper.New(),
	}
}

func (b *Builder) readConfigDir(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return errors.Wrap(err, "reading config files from dir")
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name())
		if len(ext) != 0 {
			ext = ext[1:]
		}
		if util.StringInSlice(ext, viper.SupportedExts) {
			file := filepath.Join(dir, f.Name())
			log.Debugf("merging config file %s", file)
			b.V.SetConfigFile(file)
			err = b.V.MergeInConfig()
			if err != nil {
				return errors.Wrapf(err, "merging config file '%s'", filepath.Join(dir, f.Name()))
			}
		}
	}
	return nil
}

func (b *Builder) readConfigs() error {
	for _, config := range b.configs {
		finfo, err := os.Stat(config)
		if err != nil {
			return errors.Wrap(err, "config file info")
		}

		if finfo.Mode().IsDir() {
			if err := b.readConfigDir(config); err != nil {
				return err
			}
			continue
		}

		log.Debugf("merging config file %s", config)
		b.V.SetConfigFile(config)
		err = b.V.MergeInConfig()
		if err != nil {
			return errors.Wrapf(err, "merging config file '%s'", config)
		}
	}
	return nil
}
