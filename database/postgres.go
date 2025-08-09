package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/luckydevil2007/audionotes/config"
)

func OpenPG() *sql.DB {
	conf := config.NewConfig()
	dsn := config.FormatDSN(conf.PostgresConfig)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Failed to open postgress:", err)
		return nil
	}

	return db
}
