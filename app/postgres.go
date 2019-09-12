package app

import (
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var pgOnce sync.Once
var pgClient *gorm.DB

func GetPostgres() *gorm.DB {
	pgOnce.Do(func() {
		var err error
		pgClient, err = gorm.Open("postgres", GetInstance().Config.PostgresDsn)
		if IsDebug() {
			pgClient = pgClient.Debug()
		}
		if err != nil {
			panic(err)
		}
	})
	return pgClient
}
