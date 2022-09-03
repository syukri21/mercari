package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	commonModel "github.com/syukri21/mercari/common/model"
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

func InitConfig() (model.Config, error) {
	config := model.Config{}

	// APP
	config.App = model.AppConfig{
		Env:               viper.GetString(`APP_ENV`),
		GRPCPort:          viper.GetInt(`APP_GRPC_AREA_AREA_PORT`),
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

	// DB POSTGRESSQL MASTER
	config.Postgre = commonModel.PostgreSqlConfig{
		Host:            viper.GetString("DB_POSTGRESQL_HOST"),
		Port:            viper.GetInt("DB_POSTGRESQL_PORT"),
		Username:        viper.GetString("DB_POSTGRESQL_USERNAME"),
		Password:        viper.GetString("DB_POSTGRESQL_PASSWORD"),
		MaxOpenConns:    viper.GetInt("DB_POSTGRESQL_MAX_OPEN_CONNS"),
		MaxIdleConns:    viper.GetInt("DB_POSTGRESQL_MAX_IDLE_CONNS"),
		ConnMaxLifetime: viper.GetInt("DB_POSTGRESQL_CONN_MAX_LIFETIME"),
		DBname:          viper.GetString("DB_POSTGRESQL_DBNAME"),
	}

	if config.Postgre.Host == "" ||
		config.Postgre.Port == 0 ||
		config.Postgre.Username == "" ||
		config.Postgre.Password == "" ||
		config.Postgre.MaxOpenConns == 0 ||
		config.Postgre.MaxIdleConns == 0 {
		return config, fmt.Errorf("[CONFIG][Critical] Please check section DB Postgre: %+v", config.Postgre)
	}

	return config, nil
}
