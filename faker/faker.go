package faker

import (
	log "github.com/Sirupsen/logrus"
)

const maxFakes = 2
const clicksInterval = 10000

type FakeCreator struct {
	statRepo statRepo
	fakeRepo fakeRepo
}

func NewFakeCreator() FakeCreator {
	sr := newChStatRepo()
	fr := newPgFakeRepo()
	fc := FakeCreator{
		&sr,
		&fr,
	}

	return fc
}

func (fc FakeCreator) Process() {
	stats := fc.statRepo.findSmartStats()

	l := len(stats)
	log.Infof("stats found for analyse len: %d offers", l)
	var stat stat
	for i := range stats {
		stat = stats[i]
		log.Debugf("%d/%d", i+1, l)
		log.Debug(stat)

		if ok := needFake(stat); !ok {
			continue
		}
		click, err := fc.statRepo.findAnyClick(stat.offerId, stat.publisherId)
		if err != nil {
			log.Error("click not found", err)
			continue
		}
		fake := newSmartFake(click)
		log.Info("created fake for  offerId:", fake.OfferId, "pubId:", fake.AffiliateId, "stat:", stat)

		fc.fakeRepo.save(fake)

	}
}

func needFake(stat stat) (rsl bool) {
	if stat.convs >= maxFakes {
		log.Debug("stat.convs:", stat.convs, ">=", maxFakes, ": has enough convs")
		return
	}
	if stat.clicks >= (stat.convs+1)*clicksInterval {
		log.Debug("stat.clicks:", stat.clicks, ">=", (stat.convs+1)*clicksInterval, ":(stat.convs+1)*clicksInterval")
		rsl = true
		return
	}
	return rsl
}
