package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kr/pretty"
	"tracker-smartoffer/app"
	"tracker-smartoffer/models"
)

// Calc category, weigh of real offers included into smartoffer.

type Stat struct {
	OfferId int
	Clicks  int
	Convs   int
	Revenue float32
	Weight  float32
}

var count struct {
	new  int
	good int
	wait int
}

func Rang() {
	log.Info("cmd: RANG")

	var offers []models.Offer
	offerIds := findOfferIds()
	stats := findStat(offerIds)
	var stat Stat

	log.Info("stats len:", len(stats))
	log.Info("calculating rang...")
	for _, stat = range stats {
		category := getCat(stat)
		o := models.Offer{Id: stat.OfferId, Weight: stat.Weight, Category: category}
		offers = append(offers, o)
	}

	var o models.Offer
	for _, o = range offers {
		log.Debug(o)
		saveRang(o)
	}
	log.Info(pretty.Sprint(count))
}

func saveRang(offer models.Offer) {
	db := app.GetPostgres().Exec("update smartoffer set weight=?, category=? where offer_id=?", offer.Weight, offer.Category, offer.Id)
	if db.Error != nil {
		log.Error(db.Error)
	}
}

func getCat(stat Stat) string {
	if stat.Convs > 0 {
		count.good++
		return "good"
	}
	if stat.Convs == 0 {
		if stat.Clicks > 10000 {
			count.wait++
			return "wait"
		}
	}
	count.new++
	return "new"
}

func findOfferIds() []int {
	var offers []models.Offer
	var offerIds []int
	db := app.GetPostgres().Raw(`select offer_id id from smartoffer`).Find(&offers)
	if db.Error != nil {
		log.Error(db.Error)
	}
	var offer models.Offer
	for _, offer = range offers {
		offerIds = append(offerIds, offer.Id)
	}
	log.Debug("offer_id from smartoffer:", offerIds)

	return offerIds
}
func findStat(offerIds []int) []Stat {
	rows, err := app.GetClickhouse().Query(`
select offer_id,
       sum(clicks)                 clicks,
       sum(convs)                  convs,
       sum(revenue)                revenue,
       (revenue / clicks) * 100000 weight
from (
      select offer_id,
             sum(clicks)  clicks,
             sum(convs)   convs,
             sum(revenue) revenue
      from (
             select offer_id, clicks
             from (
                   select offer_id, sum(clicks) clicks
                   from tracker.clicks_local_fast
                   group by offer_id
                    )
--              where clicks >= 10000
             ) any
             left join (
        select offer_id,count() convs,sum(revenue) revenue
        from tracker.conversions_local
        group by offer_id
        ) using offer_id
      group by offer_id
       )
where offer_id in (?)
group by offer_id
         -- having convs > 0
order by offer_id
`, offerIds)

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	if err != nil {
		log.Error(err)
	}

	var stats []Stat

	for rows.Next() {
		var (
			offerId int
			clicks  int
			convs   int
			revenue float32
			weight  float32
		)
		if err := rows.Scan(&offerId, &clicks, &convs, &revenue, &weight); err != nil {
			log.Error(err)
		}
		stat := Stat{OfferId: offerId, Clicks: clicks, Convs: convs, Revenue: revenue, Weight: weight}
		log.Debug(stat)
		stats = append(stats, stat)
	}
	return stats

}
