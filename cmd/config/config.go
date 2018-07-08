package config

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Anchor struct {
	Provider   string   `json:"provider"`
	CommonName string   `json:"cn"`
	AltNames   []string `json:"altNames"`
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
