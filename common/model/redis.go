package model

// RedisConfig ...
type RedisConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	MaxIdle     int    `json:"max_idle"`
	IdleTimeout int    `json:"idle_timeout"`
	MaxActive   int    `json:"max_active"`
}
