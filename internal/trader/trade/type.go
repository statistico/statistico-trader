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
