package provider

import (
	"fmt"
)

var (
	providers = make(map[string]Provider)
)

func Fetch(name string) (Provider, error) {
	if p, exists := providers[name]; exists {
		return p, nil
	}
	return nil, fmt.Errorf("no provider registered for '%s'", name)
}

func Store(name string, p Provider) error {
	if _, exists := providers[name]; exists {
		return fmt.Errorf("provider for name '%s' already exists", name)
	}

	providers[name] = p
	return nil
}
