package model

import "github.com/syukri21/mercari/common/model"

type Config struct {
	// APP
	App AppConfig `json:"app"`

	//	Redis
	Redis model.RedisConfig `json:"redis"`

	// PostgreSqlConfig
	Postgre model.PostgreSqlConfig `json:"postgre"`
}

// AppConfig ...
type AppConfig struct {
	Env               string `json:"env"`
	GRPCPort          int    `json:"grpc_port"`
	Timeout           int    `json:"timeout"`
	URL               string `json:"url"`
	ActiveJWTCacheTTL int64  `json:"active_jwt_cache_ttl"`
	JWTKey            string `json:"JWTKey"`
}
