package provider

import (
	"fmt"
)

type Store struct {
	providers map[string]Provider
}

var (
	defaultStore *Store
)

func DefaultStore() *Store {
	return defaultStore
}

func (s *Store) Fetch(name string) (Provider, bool) {
	if p, exists := s.providers[name]; exists {
		return p, true
	}
	return nil, false
}

func init() {
	defaultStore = NewStore()
}

func NewStore() *Store {
	return &Store{
		providers: make(map[string]Provider),
	}
}

func (s *Store) Remove(name string) {
	delete(s.providers, name)
}

func (s *Store) Store(name string, p Provider) error {
	if _, exists := s.providers[name]; exists {
		return fmt.Errorf("provider for name '%s' already exists", name)
	}

	s.providers[name] = p
	return nil
}
