package faker

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"github.com/kr/pretty"
	"tracker-smartoffer/app"
)

type fakeRepo interface {
	save(f fake)
}

type pgFakeRepo struct {
	db *gorm.DB
}

func newPgFakeRepo() pgFakeRepo {
	return pgFakeRepo{app.GetPostgres()}
}

func (r *pgFakeRepo) save(f fake) {
	rp, err := newRawpost(f)
	if err != nil {
		err = errors.New(fmt.Sprint("Error create new rawpost", err))
		return
	}
	db := app.GetPostgres().Create(&rp).Scan(&rp)
	if db.Error != nil {
		logrus.Println("Fault save fakecreator", db.Error)
		return
	}
	id := rp.Id
	pretty.Log("rawpost id:", id)
	return
}
