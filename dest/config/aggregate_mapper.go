package config

import (
	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/dest"
)

type Aggregate struct {
	Dests []Dest
}

type DestMaker func(map[string]interface{}) (dest.Dest, error)

type Mapper func(v interface{}) (dest.Dest, error)

type MapperFetcher func(id string) (DestMaker, error)

const (
	DestOpts = "destOpts"
)

func GetAggregateMapper(fetcher MapperFetcher) Mapper {
	return func(v interface{}) (dest.Dest, error) {
		conf, ok := v.(*Aggregate)
		if !ok {
			return nil, errors.New("parsing aggregate dest config unknown config type")
		}

		dests := []dest.Dest{}
		for _, aggDest := range conf.Dests {
			if aggDest.Dest == "" {
				return nil, errors.New("aggregate dest missing 'dest' key")
			}

			maker, err := fetcher(aggDest.Dest)
			if err != nil {
				return nil, errors.Wrapf(err, "fetching dest maker for %s", aggDest.Dest)
			}

			dest, err := maker(aggDest.DestOpts)
			if err != nil {
				return nil, errors.Wrapf(err, "calling maker for %s", aggDest.Dest)
			}
			dests = append(dests, dest)
		}
		return dest.NewAggregate(dests...), nil
	}
}

func NewAggregate() interface{} {
	return &Aggregate{}
}
