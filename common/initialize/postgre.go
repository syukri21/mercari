package initialize

import (
	"fmt"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/syukri21/mercari/common/model"
)

// NewPostSqlServer ...
func NewPostSqlServer(config model.PostgreSqlConfig) *sqlx.DB {
	var err error
	var postgresSqlDriverName = "postgres"

	dbURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.Username, config.Password, config.DBname)

	dbPool, err := sqlx.Connect(postgresSqlDriverName, dbURI)
	if err != nil {
		log.Panic("postgreSql.Error: ", err)
	}

	if err = dbPool.Ping(); err != nil {
		panic(err)
	}

	dbPool.SetMaxOpenConns(config.MaxOpenConns)                                    // The default is 0 (unlimited)
	dbPool.SetMaxIdleConns(config.MaxIdleConns)                                    // defaultMaxIdleConns = 2
	dbPool.SetConnMaxLifetime(time.Second * time.Duration(config.ConnMaxLifetime)) // 0, connections are reused forever.

	return dbPool
}
