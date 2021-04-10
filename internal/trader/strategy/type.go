package strategy

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"time"
)

const (
	Active = "ACTIVE"
	PercentageStakingPlan = "PERCENTAGE"
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

type StakingPlan struct {
	Name   string  `json:"name"`
	Number float32 `json:"value"`
}

func (s StakingPlan) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StakingPlan) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &s)
}
