package model

type Config struct {
	// APP
	App AppConfig `json:"app"`
}

// AppConfig ...
type AppConfig struct {
	Env               string `json:"env"`
	GRPCPort          int    `json:"grpc_port"`
	Timeout           int    `json:"timeout"`
	URL               string `json:"url"`
	ActiveJWTCacheTTL int64  `json:"active_jwt_cache_ttl"`
}
