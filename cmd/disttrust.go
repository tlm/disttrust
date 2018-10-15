package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/tlmiller/disttrust/cmd/config"
	"github.com/tlmiller/disttrust/conductor"
	"github.com/tlmiller/disttrust/provider"
	"github.com/tlmiller/disttrust/server"
	"github.com/tlmiller/disttrust/server/healthz"
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
	var apiServ *server.ApiServer
	apiServSetup := sync.Once{}
	healthApi := healthz.New()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGHUP)

	for {
		userConfig, err := config.FromFiles(configFiles...)
		if err != nil {
			log.Fatalf("making config: %v", err)
		}

		apiServSetup.Do(func() {
			if userConfig.Api.Address != "" {
				apiServ = server.NewApiServer(userConfig.Api.Address)
				healthApi.InstallHandler(apiServ.Mux)
				go apiServ.Serve()
			}
		})

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
		mstatuses := manager.AddMembers(members...)
		checks := make([]healthz.Checker, len(mstatuses))
		for i, ms := range mstatuses {
			checks[i] = ms
		}
		healthApi.SetChecks(checks...)
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
	if apiServ != nil {
		apiServ.Stop()
	}
}
