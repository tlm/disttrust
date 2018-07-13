package config

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Action struct {
	Command []string `json:"command"`
}

type Anchor struct {
	Action      Action          `json:"action"`
	AltNames    []string        `json:"altNames"`
	CommonName  string          `json:"cn"`
	Dest        string          `json:"dest"`
	DestOptions json.RawMessage `json:"destOpts"`
	Provider    string          `json:"provider"`
}

type Config struct {
	Providers []Provider `json:"providers"`
	Anchors   []Anchor   `json:"anchors"`
}

type Provider struct {
	Id      string          `json:"id"`
	Name    string          `json:"name"`
	Options json.RawMessage `json:"options"`
}

func New(data []byte) (*Config, error) {
	config := Config{}
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.Wrap(err, "building config from json")
	}
	return &config, nil
}
