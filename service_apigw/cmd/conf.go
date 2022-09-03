package main

import (
	_ "embed"
	"time"

	"github.com/spf13/viper"
)

//go:embed generated-lura-conf.json
var GeneratedLuraConf []byte

type appConfig struct {
	ServiceName                   string
	NewRelicLicenseKey            string
	SekuritasPartner              sekuritasPartner
	Redis                         redisConf
	CfbKeyJagoTradingRegistartion string
}

type redisConf struct {
	Address      string
	MinIdleConns int
	MaxConnAge   time.Duration
}

type sekuritasPartner struct {
	Host string
	Port string
}

func GetAppConfig(v *viper.Viper) appConfig {
	return appConfig{
		ServiceName:        v.GetString("VENDOR_NEW_RELIC_SERVICE_NAME"),
		NewRelicLicenseKey: v.GetString("LICENSE_KEY"),
		SekuritasPartner: sekuritasPartner{
			Host: v.GetString("sekuritas-partner-service-host"),
			Port: v.GetString("sekuritas-partner-service-port"),
		},
		Redis: redisConf{
			Address:      v.GetString("redis.address"),
			MinIdleConns: v.GetInt("redis.min-idle-conns"),
			MaxConnAge:   v.GetDuration("redis.max-conn-age"),
		},
		CfbKeyJagoTradingRegistartion: v.GetString("cfb-key-jago-trading-registartion"),
	}
}
