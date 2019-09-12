package app

import (
	"sync"

	"database/sql"

	_ "github.com/kshvakov/clickhouse"
)

var chOnce sync.Once
var chClient *sql.DB

func GetClickhouse() *sql.DB {
	chOnce.Do(func() {

		var err error
		chClient, err = sql.Open("clickhouse", GetInstance().Config.ClickhouseDsn)
		if err != nil {
			panic(err)
		}
	})

	return chClient
}
