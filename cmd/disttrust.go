package cmd

import (
	"fmt"
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
	Use:              "disttrust",
	Short:            "disttrust is a daemon that maintains local TLS certs",
	Long:             `disttrust is a daemon that maintains local TLS certs on the system through one or more providers`,
	PersistentPreRun: preRun,
	Run:              Run,
}

func applyProviders(cnfProviders []config.Provider, store *provider.Store) error {
	for _, cnfProvider := range cnfProviders {
		if len(cnfProvider.Name) == 0 {
			return errors.New("undefined provider name")
		}

		p, err := config.ToProvider(cnfProvider.Id, cnfProvider.Options)
		if err != nil {
			return errors.Wrapf(err, "config to provider%s", cnfProvider.Name)
		}

		err = store.Store(cnfProvider.Name, p)
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

func preRun(cmd *cobra.Command, args []string) {
	setupLogging()
}

func Run(cmd *cobra.Command, args []string) {
	config, err := config.FromFiles(configFiles...)
	if err != nil {
		log.Fatalf("making config: %v", err)
	}

	err = applyProviders(config.Providers, provider.DefaultStore())
	if err != nil {
		log.Fatalf("applying providers: %v", err)
	}

	err = startAnchors(config.Anchors, provider.DefaultStore())
	if err != nil {
		log.Fatalf("starting anchors: %v", err)
	}

	manager.Watch()
}

func startAnchors(anchors []config.Anchor, store *provider.Store) error {
	for _, cnfAnchor := range anchors {
		prv, err := store.Fetch(cnfAnchor.Provider)
		if err != nil {
			return errors.Wrap(err, "getting anchor provider")
		}

		req := provider.Request{}
		req.CommonName = cnfAnchor.CommonName
		req.AltNames = cnfAnchor.AltNames

		dest, err := config.ToDest(cnfAnchor.Dest, cnfAnchor.DestOptions)
		if err != nil {
			return errors.Wrap(err, "make dest for anchor")
		}
		action, err := config.ToAction(cnfAnchor.Action)
		if err != nil {
			return errors.Wrap(err, "make action for anchor")
		}

		memHandle := conductor.DefaultLeaseHandle(dest, action)
		member := conductor.NewMember(prv, req, memHandle)
		manager.AddMember(member)
	}
	return nil
}
