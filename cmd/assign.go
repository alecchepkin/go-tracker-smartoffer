package cmd

import (
	log "github.com/Sirupsen/logrus"
	"tracker-smartoffer/app"
)

const pubSystem = 1248

//Assign offers to smartoffer by geo+package_name
func Assign() {
	log.Info("cmd: ASSIGN")

	db := app.GetPostgres().Exec(`
insert into smartoffer (offer_id, smart_offer_id)
  (select o.id offer_id, so.id smart_offer_id
   from offer so
          left join offer o on so.geo = o.geo and so.package_name = o.package_name
          left join smartoffer s on o.id = s.offer_id
   where so.is_smart
     and so.status = 'active' and o.status='active'
     and so.id != o.id
     and s.offer_id isnull)
`)

	if db.Error != nil {
		log.Error(db.Error)
	}

	// assign offers to system user that should be visible in v_publisher_feed_capped_fast
	db = app.GetPostgres().Exec(`
insert into offer_to_publisher (offer_id, publisher_id, is_active) select offer_id,?,true 
from (select DISTINCT offer_id from smartoffer) t on CONFLICT do NOTHING
`, pubSystem)

	if db.Error != nil {
		log.Error(db.Error)
	}

}
