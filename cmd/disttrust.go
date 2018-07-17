package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/tlmiller/disttrust/cmd/config"
	"github.com/tlmiller/disttrust/conductor"
	"github.com/tlmiller/disttrust/provider"
)

var (
	configFiles []string
	manager     *conductor.Conductor
)

var disttrustCmd = &cobra.Command{
	Use:   "disttrust",
	Short: "disttrust is a daemon that maintains local TLS certs",
	Long:  `disttrust is a daemon that maintains local TLS certs on the system through one or more providers`,
	Run:   Run,
}

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

func Execute() {
	if err := disttrustCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	manager = conductor.NewConductor()
	disttrustCmd.Flags().StringSliceVarP(&configFiles, "config", "c",
		[]string{}, "Config file(s)")
	disttrustCmd.MarkFlagRequired("config")
}

func Run(cmd *cobra.Command, args []string) {
	for _, cFile := range configFiles {
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

	manager.Watch()
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

		dest, err := config.MakeDest(cnfAnchor.Dest, cnfAnchor.DestOptions)
		if err != nil {
			return errors.Wrap(err, "make dest for anchor")
		}
		action, err := config.MakeAction(cnfAnchor.Action)
		if err != nil {
			return errors.Wrap(err, "make action for anchor")
		}

		memHandle := conductor.DefaultLeaseHandle(dest, action)
		member := conductor.NewMember(prv, req, memHandle)
		manager.AddMember(member)
	}
	return nil
}
