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

	AwayTeam = "AWAY_TEAM"
	HomeTeam = "HOME_TEAM"

	Gte     = "GTE"
	Lte     = "LTE"
	Average = "AVERAGE"

	Continuous = "CONTINUOUS"
	Total      = "TOTAL"

	ActionFor     = "FOR"
	ActionAgainst = "AGAINST"

	Goals       = "GOALS"
	ShotsOnGoal = "SHOTS_ON_GOAL"

	Win      = "WIN"
	WinDraw  = "WIN_DRAW"
	WinLose  = "WIN_LOSE"
	Lose     = "LOSE"
	LoseDraw = "LOSE_DRAW"

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

type Result string

type Fixture struct {
	ID         uint64
	HomeTeamID uint64
	AwayTeamID uint64
	Date       time.Time
	SeasonID   uint64
}

type MatcherQuery struct {
	EventID       uint64
	ResultFilters []*ResultFilter
	StatFilters   []*StatFilter
}
