package config

import (
	"github.com/spf13/viper"

	"github.com/tlmiller/disttrust/server"
)

const (
	APIAddress = "api.address"
)

func GetAPIServer(v *viper.Viper) *server.APIServer {
	return server.NewAPIServer(v.GetString(APIAddress))
}
