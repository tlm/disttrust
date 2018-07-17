package cmd

import (
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	logJson bool = false
)

func init() {
	disttrustCmd.Flags().BoolVar(&logJson, "log-json", false, "enables json logging")
}

func setupLogging() {
	if logJson {
		log.SetFormatter(&log.JSONFormatter{
			DisableTimestamp: false,
			TimestampFormat:  time.RFC3339,
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			DisableColors:    false,
			DisableSorting:   true,
			DisableTimestamp: false,
			FullTimestamp:    true,
			TimestampFormat:  time.RFC3339,
		})
	}
}
