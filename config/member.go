package config

import (
	"fmt"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/tlmiller/disttrust/conductor"
	"github.com/tlmiller/disttrust/provider"
)

type AltNamesConfig struct {
	DNSNames    []string
	IPAddresses []string
}

type AnchorConfig struct {
	Action     ActionConfig
	AltNames   AltNamesConfig
	CommonName string
	CN         string
	Dest       string
	DestOpts   map[string]interface{}
	Name       string
	Provider   string
}

type ProviderFetcher func(string) (provider.Provider, bool)

const (
	Anchors = "anchors"
)

func GetMembersWithProviderStore(v *viper.Viper,
	store *provider.Store) ([]conductor.Member, error) {
	return GetMembers(v, ProviderFetcher(store.Fetch))
}

func GetMembers(v *viper.Viper, pFetcher ProviderFetcher) ([]conductor.Member, error) {
	anchors := []AnchorConfig{}
	members := []conductor.Member{}
	v.UnmarshalKey(Anchors, &anchors)

	for _, anchor := range anchors {
		aProvider, exists := pFetcher(anchor.Provider)
		if !exists {
			return members, fmt.Errorf("no provider found for %s", anchor.Provider)
		}

		req := provider.Request{
			AltNames: provider.AltNames{
				DNSNames:    anchor.AltNames.DNSNames,
				IPAddresses: anchor.AltNames.IPAddresses,
			},
		}

		if anchor.CN != "" {
			log.Warn("anchor cn field is deprecated in favour of commonName")
			req.CommonName = anchor.CN
		} else {
			req.CommonName = anchor.CommonName
		}

		dest, err := GetDest(&DestConfig{
			Dest:     anchor.Dest,
			DestOpts: anchor.DestOpts,
		})
		if err != nil {
			return members, errors.Wrapf(err, "failed getting dest for member %s",
				anchor.Name)
		}

		action, err := GetAction(&anchor.Action)
		if err != nil {
			return members, errors.Wrapf(err, "failed getting anchor for member %s",
				anchor.Name)
		}

		members = append(members, conductor.NewMember(anchor.Name, aProvider,
			req, conductor.DefaultLeaseHandle(dest, action)))
	}
	return members, nil
}
