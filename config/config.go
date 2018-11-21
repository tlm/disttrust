package config

import (
	"github.com/spf13/viper"
)

func SetDefaults(v *viper.Viper) {
	v.SetDefault(APIAddress, "localhost:1122")
	v.SetDefault(LogJSON, false)
	v.SetDefault(LogLevel, "info")
}
