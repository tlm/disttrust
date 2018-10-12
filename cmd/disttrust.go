package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/tlmiller/disttrust/cmd/config"
	"github.com/tlmiller/disttrust/conductor"
	"github.com/tlmiller/disttrust/provider"
	_ "github.com/tlmiller/disttrust/server"
)

var (
	configFiles []string
)

var disttrustCmd = &cobra.Command{
	Use:              "disttrust",
	Short:            "disttrust is a daemon that maintains local TLS certs",
	Long:             `disttrust is a daemon that maintains local TLS certs on the system through one or more providers`,
	PersistentPreRun: preRun,
	Run:              Run,
}

func Execute() {
	if err := disttrustCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	disttrustCmd.Flags().StringSliceVarP(&configFiles, "config", "c",
		[]string{}, "Config file(s)")
	disttrustCmd.MarkFlagRequired("config")
}

func preRun(cmd *cobra.Command, args []string) {
	setupLogging()
}

func Run(cmd *cobra.Command, args []string) {
	providers := provider.DefaultStore()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGHUP)

	for {
		userConfig, err := config.FromFiles(configFiles...)
		if err != nil {
			log.Fatalf("making config: %v", err)
		}

		log.Debug("building provider store from config")
		providers, err = config.ToProviderStore(userConfig.Providers, providers)
		if err != nil {
			log.Fatalf("building providers: %v", err)
		}

		log.Debug("building anchors from config")
		members, err := config.AnchorsToMembers(userConfig.Anchors, providers.Fetch)
		if err != nil {
			log.Fatalf("building providers: %v", err)
		}

		manager := conductor.NewConductor()
		_ = manager.AddMembers(members...)

		//var apiServ *server.ApiServer
		//if userConfig.Api.Address != "" {
		//	apiServ = server.NewApiServer(userConfig.Api.Address)
		//	for _, s := range mstatuses {
		//		apiServ.AddHealthzChecks(s)
		//	}
		//	go apiServ.Serve()
		//}

		manager.Play()

		sig := <-sigCh
		log.Infof("recieved signal %s", sig)

		if sysSig, ok := sig.(syscall.Signal); ok && sysSig == syscall.SIGHUP {
			manager.Stop()
			log.Info("reloading config")
		} else {
			break
		}
	}
}
