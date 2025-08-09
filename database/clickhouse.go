package database

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/luckydevil2007/audionotes/config"
)

func OpenClickhouse() *sql.DB {
	conf := config.NewConfig().ClickhouseConfig
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{conf.Host + ":" + strconv.Itoa(conf.Port)},

		//	DialTimeout: 5 * time.Second,
	})

	/*
		conf := config.NewConfig()
		dsn := config.FormatDSN(conf)
		db, err := sql.Open("clickhouse", dsn)*/
	if db == nil {
		fmt.Println("Failed to open clickhouse:")
		return nil
	}

	return db
}
