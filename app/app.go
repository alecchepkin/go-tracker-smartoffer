package app

import (
	"sync"

	"github.com/caarlos0/env"
)

type config struct {
	PostgresDsn   string `env:"POSTGRES_DSN" envDefault:"host=localhost port=5432 user=postgres dbname=postgres sslmode=disable"`
	ClickhouseDsn string `env:"CLICKHOUSE_DSN" envDefault:"http://localhost:8123"`
	LogLevel      string `env:"LOG_LEVEL" envDefault:"debug"`
}

type App struct {
	Config config
}

var once sync.Once
var instance *App

func GetInstance() *App {
	once.Do(func() {
		instance = &App{}
		err := env.Parse(&instance.Config)
		if err != nil {
			panic(err)
		}
	})
	return instance
}

func IsDebug() bool {
	lvl := GetInstance().Config.LogLevel

	if lvl == "debug" {
		return true
	}

	return false
}
