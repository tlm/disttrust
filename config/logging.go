package config

import (
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

const (
	LogJSON  string = "logging.json"
	LogLevel string = "logging.level"
)

func SetLogging(v *viper.Viper) error {
	if v.GetBool(LogJSON) {
		logrus.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: false,
			TimestampFormat:  time.RFC3339,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors:    false,
			DisableSorting:   true,
			DisableTimestamp: false,
			FullTimestamp:    true,
			TimestampFormat:  time.RFC3339,
		})
	}
	lvl, err := logrus.ParseLevel(v.GetString(LogLevel))
	if err != nil {
		return errors.Wrap(err, "parsing config log level")
	}
	logrus.SetLevel(lvl)
	return nil
}
