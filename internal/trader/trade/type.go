package trade

import (
	"github.com/google/uuid"
	"time"
)

const (
	InPlay = "IN_PLAY"
)

type Trade struct {
	ID          uuid.UUID `json:"id"`
	StrategyID  uuid.UUID `json:"strategyId"`
	Exchange    string    `json:"exchange"`
	ExchangeRef string    `json:"exchangeRef"`
	Market      string    `json:"market"`
	Runner      string    `json:"runner"`
	Price       float32   `json:"price"`
	Stake       float32   `json:"stake"`
	EventID     uint64    `json:"eventId"`
	EventDate   time.Time `json:"eventDate"`
	Side        string    `json:"side"`
	Result      string    `json:"result"`
	Timestamp   time.Time `json:"timestamp"`
}

type Ticket struct {
	MarketID      string    `json:"marketId"`
	MarketName    string    `json:"marketName"`
	RunnerID      uint64    `json:"runnerId"`
	RunnerName    string    `json:"runnerName"`
	EventID       uint64    `json:"eventId"`
	CompetitionID uint64    `json:"competitionId"`
	SeasonID      uint64    `json:"seasonId"`
	EventDate     time.Time `json:"date"`
	Exchange      string    `json:"exchange"`
	Price         TicketPrice     `json:"price"`
}

type TicketPrice struct {
	Value     float32   `json:"price"`
	Size      float32   `json:"size"`
	Side      string    `json:"side"`
}
