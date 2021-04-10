package queue

import "time"

type MarketQueue interface {
	ReceiveMarkets() <-chan *EventMarket
}

type EventMarket struct {
	ID            string    `json:"id"`
	EventID       uint64    `json:"eventId"`
	Name          string    `json:"name"`
	CompetitionID uint64    `json:"competitionId"`
	SeasonID      uint64    `json:"seasonId"`
	EventDate     time.Time `json:"date"`
	Exchange      string    `json:"exchange"`
	Runners       []*Runner `json:"runners"`
	Timestamp     int64     `json:"timestamp"`
}

type Runner struct {
	ID     uint64      `json:"id"`
	Name   string      `json:"name"`
	Sort   int8        `json:"sort"`
	BackPrices []PriceSize `json:"backPrices"`
	LayPrices  []PriceSize `json:"layPrices"`
}

type PriceSize struct {
	Price float32 `json:"price"`
	Size  float32 `json:"size"`
}
