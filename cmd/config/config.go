package config

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Anchor struct {
	Action      Action          `json:"action"`
	AltNames    []string        `json:"altNames"`
	CommonName  string          `json:"cn"`
	Dest        string          `json:"dest"`
	DestOptions json.RawMessage `json:"destOpts"`
	Name        string          `json:"name"`
	Provider    string          `json:"provider"`
}

type Config struct {
	Api       Api        `json:"api"`
	Providers []Provider `json:"providers"`
	Anchors   []Anchor   `json:"anchors"`
}

func DefaultConfig() *Config {
	return &Config{}
}

func mergeConfigs(dst *Config, src *Config) *Config {
	if len(src.Api.Address) != 0 {
		dst.Api = src.Api
	}
	dst.Providers = append(dst.Providers, src.Providers...)
	dst.Anchors = append(dst.Anchors, src.Anchors...)
	return dst
}

func New(data []byte) (*Config, error) {
	config := Config{}
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.Wrap(err, "building config from json")
	}
	return &config, nil
}
