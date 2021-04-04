package market

import (
	"github.com/statistico/statistico-strategy/internal/trader"
	"time"
)

const (
	Away = "Away"
	Draw = "Draw"
	Home = "Home"

	MatchOdds   = "MATCH_ODDS"
	OverUnder05 = "OVER_UNDER_05"
	OverUnder15 = "OVER_UNDER_15"
	OverUnder25 = "OVER_UNDER_25"
	OverUnder35 = "OVER_UNDER_35"
	OverUnder45 = "OVER_UNDER_45"

	Fail    = "FAIL"
	Success = "SUCCESS"

	Back = "BACK"
	Lay  = "LAY"

	Over  = "Over"
	Under = "Under"
)

type Trade struct {
	MarketName    string    `json:"marketName"`
	RunnerName    string    `json:"runnerName"`
	RunnerPrice   float32   `json:"runnerPrice"`
	EventId       uint64    `json:"eventId"`
	CompetitionId uint64    `json:"competitionId"`
	SeasonId      uint64    `json:"seasonId"`
	EventDate     time.Time `json:"eventDate"`
	Side          string    `json:"side"`
	Exchange      string    `json:"exchange"`
	Result        *string   `json:"result,omitempty"`
}

type Query struct {
	MarketName    string
	RunnerName    string
	RunnerPrice   float32
	EventId       uint64
	CompetitionId uint64
	SeasonId      uint64
	EventDate     time.Time
	Side          string
	Exchange      string
	ResultFilters []*trader.ResultFilter
	StatFilters   []*trader.StatFilter
}

type MatcherQuery struct {
	EventID       uint64
	ResultFilters []*trader.ResultFilter
	StatFilters   []*trader.StatFilter
}