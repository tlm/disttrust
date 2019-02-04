package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/tlmiller/disttrust/conductor"
	"github.com/tlmiller/disttrust/config"
	"github.com/tlmiller/disttrust/provider"
	"github.com/tlmiller/disttrust/server"
	"github.com/tlmiller/disttrust/server/healthz"
)

const (
	flagConfig   = "config"
	flagLogJSON  = "log-json"
	flagLogLevel = "log-level"
)

var (
	configFiles []string
)

var disttrustCmd = &cobra.Command{
	Use:   "disttrust",
	Short: "disttrust is a daemon that maintains local TLS certs",
	Long:  `disttrust is a daemon that maintains local TLS certs on the system through one or more providers`,
	Run:   Run,
}

func Execute() {
	if err := disttrustCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	disttrustCmd.Flags().StringSliceVarP(&configFiles, flagConfig, "c",
		[]string{}, "One or more config files and or config directories")
	disttrustCmd.Flags().Bool(flagLogJSON, false, "log messages in json format")
	disttrustCmd.Flags().String(flagLogLevel, "",
		"level to log messages at, one of panic, fatal, error, warn, info, debug or trace")
	disttrustCmd.MarkFlagRequired(flagConfig)
}

func Run(cmd *cobra.Command, args []string) {
	providers := provider.DefaultStore()
	var apiServ *server.APIServer
	apiServSetup := sync.Once{}
	healthApi := healthz.New()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGHUP)

	for {
		confB := config.NewBuilder(configFiles)
		registerConfigFlags(cmd.Flags(), confB.V)
		err := confB.BuildAndValidate()
		if err != nil {
			log.Fatalf("making config: %v", err)
		}

		if err := config.SetLogging(confB.V); err != nil {
			log.Fatalf("configuring logging: %v\n", err)
		}
		log.Debug("finished setting up logger")

		apiServSetup.Do(func() {
			log.Debug("building api server from config")
			apiServ = config.GetAPIServer(confB.V)
			healthApi.InstallHandler(apiServ.Mux)
			log.Debugf("starting api server on %s", apiServ.Server.Addr)
			go apiServ.Serve()
		})

		log.Debug("building provider store from config")
		providers, err = config.GetProviderStore(confB.V, providers)
		if err != nil {
			log.Fatalf("building providers: %v", err)
		}

		log.Debug("initialising providers")
		for _, provider := range providers {
			log := log.WithFields(log.Fields{
				"provider_name": provider.Name(),
			})
			log.Info("initialising provider")
			if err := provider.Initialise(); err != nil {
				log.Fatalf("initialising provider: %v", err)
			}
		}

		log.Debug("building members from config")
		members, err := config.GetMembersWithProviderStore(confB.V, providers)
		if err != nil {
			log.Fatalf("building members: %v", err)
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

func registerConfigFlags(flags *pflag.FlagSet, v *viper.Viper) {
	v.BindPFlag(config.LogJSON, flags.Lookup(flagLogJSON))
	v.BindPFlag(config.LogLevel, flags.Lookup(flagLogLevel))
}
