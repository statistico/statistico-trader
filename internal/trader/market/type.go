package market

import "time"

type Runner struct {
	MarketID      string    `json:"marketId"`
	MarketName    string    `json:"marketName"`
	RunnerID      uint64    `json:"runnerId"`
	RunnerName    string    `json:"runnerName"`
	EventID       uint64    `json:"eventId"`
	CompetitionID uint64    `json:"competitionId"`
	SeasonID      uint64    `json:"seasonId"`
	EventDate     time.Time `json:"date"`
	Exchange      string    `json:"exchange"`
	Price         Price     `json:"price"`
}

type Price struct {
	Value     float32   `json:"price"`
	Size      float32   `json:"size"`
	Side      string    `json:"side"`
	Timestamp time.Time `json:"timestamp"`
}
