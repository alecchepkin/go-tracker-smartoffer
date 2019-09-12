package cmd

import (
	log "github.com/Sirupsen/logrus"
	"tracker-smartoffer/faker"
)

func Fake() {
	log.Infoln("cmd: FAKE")
	faker.NewFakeCreator().Process()
}
