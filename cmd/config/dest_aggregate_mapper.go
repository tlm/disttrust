package config

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/dest"
)

func destAggregateMapper(opts json.RawMessage) (dest.Dest, error) {
	uopts := []Dest{}
	err := json.Unmarshal(opts, &uopts)
	if err != nil {
		return nil, errors.Wrap(err, "parsing dest aggregate json")
	}

	dests := []dest.Dest{}
	for _, rawDest := range uopts {
		if rawDest.Dest == "" {
			return nil, errors.New("aggregate dest missing 'dest' key")
		}
		dest, err := ToDest(rawDest.Dest, rawDest.DestOptions)
		if err != nil {
			return nil, errors.Wrap(err, "aggregate dest failed parsing")
		}
		dests = append(dests, dest)
	}
	return dest.NewAggregate(dests...), nil
}

func init() {
	err := MapDest("aggregate", destAggregateMapper)
	if err != nil {
		panic(err)
	}
}
