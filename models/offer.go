package models

type Offer struct {
	Id           int
	SmartOfferId int
	Revenue      float32
	Payout       float32
	Category     string
	Weight       float32
}

func (Offer) TableName() string {
	return "public.offer"
}
