package config

import "encoding/json"

type Provider struct {
	Id      string          `json:"id"`
	Name    string          `json:"name"`
	Options json.RawMessage `json:"options"`
}
