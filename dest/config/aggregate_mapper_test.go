package config

import (
	"testing"

	"github.com/tlmiller/disttrust/dest"
)

type test struct {
	Dest string
}

func TestAggregateMapperMissingData(t *testing.T) {
	testFetcher := func(_ string) (DestMaker, error) {
		return nil, nil
	}

	mapper := GetAggregateMapper(testFetcher)
	if _, err := mapper(NewAggregate()); err != nil {
		t.Fatalf("unexpected error calling aggregate mapper: %v", err)
	}
}

func TestAggregateMapperCallsMapper(t *testing.T) {
	called := false
	testFetcher := func(id string) (DestMaker, error) {
		if id != "test" {
			t.Fatalf("unexepected dest id wanted test got %s", id)
		}
		return func(_ map[string]interface{}) (dest.Dest, error) {
			called = true
			return nil, nil
		}, nil
	}

	mapper := GetAggregateMapper(testFetcher)
	conf := &Aggregate{}
	conf.Dests = append(conf.Dests, Dest{
		Dest: "test",
	})
	if _, err := mapper(conf); err != nil {
		t.Fatalf("unexpected error calling aggregate mapper: %v", err)
	}

	if !called {
		t.Error("aggregate dest mapper was never called")
	}
}
