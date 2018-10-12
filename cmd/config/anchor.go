package config

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/conductor"
	"github.com/tlmiller/disttrust/provider"
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

type ProviderFixer func(string) (provider.Provider, bool)

func AnchorsToMembers(anchors []Anchor, fixer ProviderFixer) ([]conductor.Member, error) {
	members := []conductor.Member{}
	for _, anchor := range anchors {
		aprovider, exists := fixer(anchor.Provider)
		if !exists {
			return nil, fmt.Errorf("no provider found for %s", anchor.Provider)
		}

		req := provider.Request{}
		req.CommonName = anchor.CommonName
		req.AltNames = anchor.AltNames

		dest, err := ToDest(anchor.Dest, anchor.DestOptions)
		if err != nil {
			return nil, errors.Wrap(err, "making dest for anchor")
		}
		action, err := ToAction(anchor.Action)
		if err != nil {
			return nil, errors.Wrap(err, "making action for anchor")
		}

		members = append(members, conductor.NewMember(anchor.Name, aprovider,
			req, conductor.DefaultLeaseHandle(dest, action)))
	}
	return members, nil
}
