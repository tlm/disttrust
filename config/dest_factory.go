package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/dest"
	destconf "github.com/tlmiller/disttrust/dest/config"
)

type DestConfigFunc func() interface{}

type DestMapper func(v interface{}) (dest.Dest, error)

type DestMapping struct {
	Config DestConfigFunc
	Mapper DestMapper
}

type DestFactory map[string]DestMapping

var (
	DefaultDestFactory = make(DestFactory)
)

func (d DestFactory) AddMapping(id string, m DestMapping) error {
	if _, exists := d[id]; exists {
		return fmt.Errorf("dest factory mapping already exists for id %s", id)
	}
	d[id] = m
	return nil
}

func init() {
	aggFetcher := func(id string) (destconf.DestMaker, error) {
		mapping, exists := DefaultDestFactory[id]
		if !exists {
			return nil, fmt.Errorf("dest maker does not exist for id %s", id)
		}
		return destconf.DestMaker(mapping.MakeDest), nil
	}
	aggregateMapping := DestMapping{
		Config: destconf.NewAggregate,
		Mapper: DestMapper(destconf.GetAggregateMapper(aggFetcher)),
	}
	if err := DefaultDestFactory.AddMapping("aggregate", aggregateMapping); err != nil {
		panic(err)
	}

	fileMapping := DestMapping{
		Config: destconf.NewFile,
		Mapper: destconf.FileMapper,
	}
	if err := DefaultDestFactory.AddMapping("file", fileMapping); err != nil {
		panic(err)
	}

	templateMapping := DestMapping{
		Config: destconf.NewTemplate,
		Mapper: destconf.TemplateMapper,
	}
	if err := DefaultDestFactory.AddMapping("template", templateMapping); err != nil {
		panic(err)
	}
}

func (d DestMapping) MakeDest(options map[string]interface{}) (dest.Dest, error) {
	conf := d.Config()
	if err := mapstructure.Decode(options, conf); err != nil {
		return nil, errors.Wrap(err, "making dest mapping config from options")
	}
	return d.Mapper(conf)
}
