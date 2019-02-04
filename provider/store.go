package provider

import (
	"fmt"
)

type Store map[string]Provider

var (
	defaultStore Store
)

func DefaultStore() Store {
	return defaultStore
}

func (s Store) Fetch(name string) (Provider, bool) {
	if p, exists := s[name]; exists {
		return p, true
	}
	return nil, false
}

func init() {
	defaultStore = NewStore()
}

func NewStore() Store {
	return Store{}
}

func (s Store) Remove(name string) {
	delete(s, name)
}

func (s Store) Store(name string, p Provider) error {
	if _, exists := s[name]; exists {
		return fmt.Errorf("provider for name '%s' already exists", name)
	}

	s[name] = p
	return nil
}
