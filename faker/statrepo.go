package faker

import (
	"database/sql"
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"github.com/kr/pretty"
	"tracker-smartoffer/app"
)

type statRepo interface {
	findSmartStats() []stat
	findAnyClick(offerId int, publisherId int) (clk click, err error)
}

func newChStatRepo() chStatRepo {
	return chStatRepo{app.GetClickhouse()}
}

type chStatRepo struct {
	db *sql.DB
}

func (r *chStatRepo) findSmartStats() []stat {
	/*	return []stat{
		stat{offerId: 1, publisherId: 1, clicks: 1, convs: 2},
		stat{offerId: 1, publisherId: 1, clicks: 1, convs: 0},
		stat{offerId: 0, publisherId: 1, clicks: 10001, convs: 0},
	}*/
	/*	rows, err := r.db.Query(`
		select offer_id, publisher_id, convs, clicks
		from (
		      select offer_id, publisher_id, sum(clicks) clicks
		      from tracker.clicks_local_fast
		      group by offer_id, publisher_id) any
		      left join (select offer_id, publisher_id, count(1) convs
		                 from tracker.conversions_local conv
		                 group by offer_id, publisher_id)
		                using offer_id, publisher_id
		where clicks > ?
		order by offer_id, publisher_id;
		`, clicksInterval)

	*/
	rows, err := app.GetClickhouse().Query(`
		select smart_offer_id, publisher_id, convs, clicks
		from (
		       select smart_offer_id, publisher_id, sum(clicks) clicks
		       from tracker.clicks_local_fast
		       where smart_offer_id > 0
		       group by smart_offer_id, publisher_id) any
		       left join (select smart_offer_id, publisher_id, count(1) convs
		                  from tracker.conversions_local conv
		                  where smart_offer_id > 0
		                  group by smart_offer_id, publisher_id)
		                 using smart_offer_id, publisher_id
		where clicks > ?
		order by smart_offer_id, publisher_id;

		`, clicksInterval)

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	if err != nil {
		logrus.Error(err)
	}

	var stats []stat

	for rows.Next() {
		var (
			offerId     int
			publisherId int
			clicks      int
			convs       int
		)
		if err := rows.Scan(&offerId, &publisherId, &clicks, &convs); err != nil {
			logrus.Error(err)
		}
		s := stat{offerId: offerId, publisherId: publisherId, clicks: clicks, convs: convs}
		pretty.Println(s)
		stats = append(stats, s)
	}
	return stats

}
func (r *chStatRepo) findAnyClick(offerId int, publisherId int) (clk click, err error) {

	data, err := r.db.Query("select id from tracker.clicks_local where click_date = today()  and smart_offer_id=? and publisher_id=?  limit 1", offerId, publisherId)

	defer func() {
		err := data.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	if err != nil {
		logrus.Error(err)
	}

	var (
		id string
	)
	if data.Next() {
		if err = data.Scan(&id); err != nil {
			logrus.Error(err)
			err = errors.New("click not found")
			return
		}
	}
	clk = click{id: uuid.MustParse(id), offerId: offerId, publisherId: publisherId}

	return
}
