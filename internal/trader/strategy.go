package trader

import (
	"github.com/google/uuid"
	"time"
)

type Strategy struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	UserID         uuid.UUID       `json:"userId"`
	MarketName     string          `json:"market"`
	RunnerName     string          `json:"runner"`
	MinOdds        *float32        `json:"minOdds"`
	MaxOdds        *float32        `json:"maxOdds"`
	CompetitionIDs []uint64        `json:"competitionIds"`
	Side           string          `json:"side"`
	Visibility     string          `json:"visibility"`
	Status         string          `json:"status"`
	StakingPlan    StakingPlan     `json:"stakingPlan"`
	ResultFilters  []*ResultFilter `json:"resultFilters"`
	StatFilters    []*StatFilter   `json:"statFilters"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type ResultFilter struct {
	Team   string `json:"team"`
	Result string `json:"result"`
	Games  uint8  `json:"games"`
	Venue  string `json:"venue"`
}

type StatFilter struct {
	Stat    string  `json:"stat"`
	Team    string  `json:"team"`
	Action  string  `json:"action"`
	Games   uint8   `json:"games"`
	Measure string  `json:"measure"`
	Metric  string  `json:"metric"`
	Value   float32 `json:"value"`
	Venue   string  `json:"venue"`
}

type StakingPlan interface {
	Identifier() string
	Stake(bank float32) float32
}

type PercentageStakingPlan struct {
	Name   string   `json:"name"`
	Value  float32  `json:"value"`
}

func (p *PercentageStakingPlan) Identifier() string {
	return p.Name
}

func (p *PercentageStakingPlan) Stake(bank float32) float32 {
	return (bank / 100) * p.Value
}

type StrategyWriter interface {
	Insert(s *Strategy) error
}
