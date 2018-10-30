package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/dest"
	"github.com/tlmiller/disttrust/file"
)

type Dest struct {
	Dest        string          `json:"dest"`
	DestOptions json.RawMessage `json:"destOpts"`
}

type DestMapper func(json.RawMessage) (dest.Dest, error)

var (
	destMappings = make(map[string]DestMapper)
)

func MapDest(id string, mapper DestMapper) error {
	if _, exists := destMappings[id]; exists {
		return fmt.Errorf("dest mapping already registered for id '%s'", id)
	}
	destMappings[id] = mapper
	return nil
}

func ToDest(id string, opts json.RawMessage) (dest.Dest, error) {
	mapper, exists := destMappings[id]
	if !exists {
		return nil, fmt.Errorf("dest mapper does not exist for id '%s'", id)
	}
	return mapper(opts)
}

func destBuilder(path, mode, gid, uid string) (file.File, error) {
	builder := file.New(path)
	if mode != "" {
		conv, err := strconv.ParseUint(mode, 8, 32)
		if err != nil {
			return file.File{}, errors.Wrap(err, "invalid mode uint")
		}
		builder.Mode = os.FileMode(conv)
	}
	builder.Gid = gid
	builder.Uid = uid
	return builder, nil
}
