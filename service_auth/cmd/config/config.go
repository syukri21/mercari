package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/syukri21/mercari/service_auth/model"
)

func init() {
	viper.SetConfigFile(`config.json`)
	// Enable VIPER to read Environment Variables
	// viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Can't find config file, err : %s", err.Error())
	}
}

func InitConfig() model.Config {
	config := model.Config{}

	// APP
	config.App = model.AppConfig{
		Env:               viper.GetString(`APP_ENV`),
		GRPCPort:          viper.GetInt(`APP_GRPC_PORT`),
		Timeout:           viper.GetInt(`APP_TIMEOUT`),
		URL:               viper.GetString(`APP_URL`),
		ActiveJWTCacheTTL: viper.GetInt64(`APP_ACTIVE_JWT_CACHE_TTL`),
	}

	// Set default to device ID cache TTL if not set
	if config.App.ActiveJWTCacheTTL == 0 {
		config.App.ActiveJWTCacheTTL = 43200
	}

	return config
}
