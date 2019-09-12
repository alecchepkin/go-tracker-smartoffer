package faker

import (
	"fmt"
	"github.com/google/uuid"
)

type fake struct {
	id            int
	Uuid          uuid.UUID `json:"uuid"`
	TransactionId string    `json:"transaction_id"`
	OfferId       int       `json:"offer_id"`
	AffiliateId   int       `json:"affiliate_id"`
}

func newSmartFake(c click) fake {
	return fake{
		Uuid:          uuid.New(),
		TransactionId: fmt.Sprintf("FK_%s", c.id),
		OfferId:       c.offerId,
		AffiliateId:   c.publisherId,
	}
}
