package main

import (
	"flag"
	"io/ioutil"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/tlmiller/disttrust/cmd/config"
	cmdflag "github.com/tlmiller/disttrust/cmd/flag"
	"github.com/tlmiller/disttrust/conductor"
	"github.com/tlmiller/disttrust/provider"
)

var (
	manager = conductor.NewConductor()
)

func buildProviders(cnfProviders []config.Provider) error {
	for _, cnfProvider := range cnfProviders {
		if len(cnfProvider.Name) == 0 {
			return errors.New("undefined provider name")
		}
		log.Debugf("building provider %s", cnfProvider.Name)

		p, err := config.MakeProvider(cnfProvider.Id, cnfProvider.Options)
		if err != nil {
			return errors.Wrapf(err, "building provider %s", cnfProvider.Name)
		}

		log.WithFields(log.Fields{
			"providerName": cnfProvider.Name,
			"providerId":   cnfProvider.Id,
		}).Info("registering provider")
		err = provider.Store(cnfProvider.Name, p)
		if err != nil {
			return errors.Wrap(err, "registering provider")
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if len(cmdflag.ConfigFiles) == 0 {
		log.Fatal("no config files provided")
	}

	for _, cFile := range cmdflag.ConfigFiles {
		log.Debugf("parsing config file %s", cFile)
		raw, err := ioutil.ReadFile(cFile)
		if err != nil {
			log.Fatalf("reading config file: %v", err)
		}

		cnf, err := config.New(raw)
		if err != nil {
			log.Fatalf("parsing config file: %v", err)
		}

		err = buildProviders(cnf.Providers)
		if err != nil {
			log.Fatalf("building providers from config file %s: %v", cFile, err)
		}

		err = startAnchors(cnf.Anchors)
		if err != nil {
			log.Fatalf("starting anchors from config file %s: %v", cFile, err)
		}
	}

	manager.Conduct()
}

func startAnchors(anchors []config.Anchor) error {
	for _, cnfAnchor := range anchors {
		prv, err := provider.Fetch(cnfAnchor.Provider)
		if err != nil {
			return errors.Wrap(err, "getting anchor provider")
		}

		req := provider.Request{}
		req.CommonName = cnfAnchor.CommonName
		req.AltNames = cnfAnchor.AltNames

		member := conductor.Member{}
		member.Provider = prv
		member.Request = req

		manager.AddMember(member)
	}
	return nil
}
