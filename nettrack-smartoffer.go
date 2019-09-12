package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"tracker-smartoffer/app"
	"tracker-smartoffer/cmd"
	"os"

	"github.com/sgreben/flagvar"
)

var (
	sub = flagvar.Enum{Choices: []string{"assign", "rang", "payout", "fake", "run"}}
)

func init() {
	lvl := app.GetInstance().Config.LogLevel
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		log.Error(err)
		ll = log.DebugLevel
	}
	log.SetLevel(ll)
	log.Info("LogLevel:", ll)
}

func main() {

	flag.Var(&sub, "sub", fmt.Sprintf("subroutine (%s)", sub.Help()))
	flag.Parse()

	fmt.Println(sub.Value)

	switch sub.Value {
	case "assign":
		cmd.Assign()
	case "fake":
		cmd.Fake()
	case "payout":
		cmd.Payout()
	case "rang":
		cmd.Rang()
	case "run":
		cmd.Assign()
		cmd.Fake()
		cmd.Payout()
		cmd.Rang()
	}

	if sub.Value != "" {
		os.Exit(0)
	}

	c := cron.New()
	err := c.AddFunc("0 */1 * * * *", func() {

		cmd.Assign()
		cmd.Fake()
		cmd.Payout()
		cmd.Rang()
		log.Infof("Cron done")
	})

	if err != nil {
		log.Fatal(err)
	}
	c.Start()
	cronjob("* */10 * * * *", "Run All Commands", func() {
		cmd.Assign()
		cmd.Fake()
		cmd.Payout()
		cmd.Rang()
	}, c)

	log.Info("Running as server")

	r := gin.Default()
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	// listen and serve on 0.0.0.0:8080
	err = r.Run()
	if err != nil {
		panic(err)
	}

}


func cronjob(spec string, name string, cmd func(), c *cron.Cron) {
	err := c.AddFunc(spec, func() {
		cmd()
	})
	if err != nil {
		panic(err)
	}
}
