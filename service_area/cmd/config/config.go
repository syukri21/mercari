package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	commonModel "github.com/syukri21/mercari/common/model"
	"github.com/syukri21/mercari/service_area/model"
)

func init() {
	viper.SetConfigFile(`config.json`)
	// Enable VIPER to read Environment Variables
	// viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Can't find config file, err : %s", err.Error())
	}
}

func InitConfig() (model.Config, error) {
	config := model.Config{}

	// APP
	config.App = model.AppConfig{
		Env:               viper.GetString(`APP_ENV`),
		GRPCPort:          viper.GetInt(`APP_GRPC_AREA_AUTH_PORT`),
		Timeout:           viper.GetInt(`APP_TIMEOUT`),
		URL:               viper.GetString(`APP_URL`),
		ActiveJWTCacheTTL: viper.GetInt64(`APP_ACTIVE_JWT_CACHE_TTL`),
	}

	// Set default to device ID cache TTL if not set
	if config.App.ActiveJWTCacheTTL == 0 {
		config.App.ActiveJWTCacheTTL = 43200
	}

	// APP
	config.Redis = commonModel.RedisConfig{
		Host:        viper.GetString(`DB_REDIS_HOST`),
		Port:        viper.GetInt(`DB_REDIS_PORT`),
		MaxIdle:     viper.GetInt(`DB_REDIS_MAX_IDLE`),
		IdleTimeout: viper.GetInt(`DB_REDIS_IDLE_TIMEOUT`),
		MaxActive:   viper.GetInt(`DB_REDIS_MAX_ACTIVE`),
	}

	if config.Redis.Host == "" ||
		config.Redis.Port == 0 ||
		config.Redis.MaxIdle == 0 ||
		config.Redis.IdleTimeout == 0 ||
		config.Redis.MaxActive == 0 {
		err := fmt.Errorf("[CONFIG][Critical] Please check section DB REDIS: %+v", config.Redis)
		return config, err
	}

	// Cron
	config.Cron = model.CronConfig{
		FillArea: viper.GetString(`CRON_FILL_AREA`),
	}

	if config.Redis.Host == "" ||
		config.Redis.Port == 0 ||
		config.Redis.MaxIdle == 0 ||
		config.Redis.IdleTimeout == 0 ||
		config.Redis.MaxActive == 0 {
		err := fmt.Errorf("[CONFIG][Critical] Please check section DB REDIS: %+v", config.Redis)
		return config, err
	}

	return config, nil
}
