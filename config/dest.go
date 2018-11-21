package config

import (
	"fmt"

	"github.com/tlmiller/disttrust/dest"
)

type DestConfig struct {
	Dest     string
	DestOpts map[string]interface{}
}

func GetDest(d *DestConfig) (dest.Dest, error) {
	if d.Dest == "" {
		return &dest.Empty{}, nil
	}

	mapping, exists := DefaultDestFactory[d.Dest]
	if !exists {
		return nil, fmt.Errorf("dest mapper does not exist for id %s", d.Dest)
	}
	return mapping.MakeDest(d.DestOpts)
}
