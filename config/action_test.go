package config

import (
	"testing"

	"github.com/tlmiller/disttrust/action"
)

func TestActionGetWithNoData(t *testing.T) {
	conf := &ActionConfig{}
	a, err := GetAction(conf)
	if err != nil {
		t.Fatalf("unexpected error for no data to action get")
	}
	_, ok := a.(*action.Empty)
	if !ok {
		t.Fatalf("unexpected type for action with no data")
	}
}
