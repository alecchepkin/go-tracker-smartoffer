package cmd

import (
	log "github.com/Sirupsen/logrus"
	"tracker-smartoffer/app"
	"tracker-smartoffer/models"
	"strconv"
)

type OfferMap map[models.Offer][]models.Offer

func Payout() {
	log.Info("cmd: PAYOUT")

	offersMap := findAllSmartOffersWithOffers()
	log.Info("stats len:", len(offersMap))

	var (
		so models.Offer
		oo []models.Offer
		i  int
	)

	for so, oo = range offersMap {
		revenue, payout := smartPayout(oo)
		so.Revenue = revenue
		so.Payout = payout
		savePayout(so)
		i++
		log.Debug(strconv.Itoa(i)+".", "payout:", payout, "set for smart offer", so.Id)
	}
}

func savePayout(offer models.Offer) {

	db := app.GetPostgres().Exec("update offer set revenue=?, default_payout=? where id=?", offer.Revenue, offer.Payout, offer.Id)
	if db.Error != nil {
		log.Error(db.Error)
	}
	db = app.GetPostgres().Exec("update offer_to_publisher set payout=? where offer_id=?", offer.Payout, offer.Id)
	if db.Error != nil {
		log.Error(db.Error)
	}
}

func smartPayout(offers []models.Offer) (float32, float32) {
	var (
		sum float32
		o   models.Offer
	)
	for _, o = range offers {
		sum += o.Revenue
	}

	sum2 := sum * 1.3

	return sum / float32(len(offers)), sum2 / float32(len(offers))
}

func findAllSmartOffersWithOffers() OfferMap {
	m := make(OfferMap)

	sql := `select *, default_payout as payout from offer o join smartoffer so on o.id=so.offer_id where so.smart_offer_id in (select id from offer where is_smart and o.status='active')  order by smart_offer_id`
	var offers []models.Offer
	db := app.GetPostgres().Raw(sql).Find(&offers)
	if db.Error != nil {
		log.Error(db.Error)
	}

	var (
		oo = make(map[int][]models.Offer)
		o  models.Offer
		ok bool
	)
	for _, o = range offers {
		if _, ok = oo[o.SmartOfferId]; !ok {
			oo[o.SmartOfferId] = []models.Offer{}
		}
		oo[o.SmartOfferId] = append(oo[o.SmartOfferId], o)
	}

	sql = `select *, default_payout as payout from offer where is_smart`
	var smartOffers []models.Offer
	db = app.GetPostgres().Raw(sql).Find(&smartOffers)
	if db.Error != nil {
		log.Error(db.Error)
	}

	var so models.Offer
	for _, so = range smartOffers {
		m[so] = oo[so.Id]
	}

	return m
}
