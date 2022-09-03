package model

// PostgreSqlConfig ...
type PostgreSqlConfig struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	DBname          string `json:"dbname"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
	DBRootCert      string `json:"db_root_cert"`
	DBCert          string `json:"db_cert"`
	DBKey           string `json:"db_key"`
}
