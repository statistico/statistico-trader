package classify

import (
	"github.com/statistico/statistico-strategy/internal/trader"
	"time"
)

const (
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
)

type Fixture struct {
	ID         uint64
	HomeTeamID uint64
	AwayTeamID uint64
	Date       time.Time
	SeasonID   uint64
}

type MatcherQuery struct {
	EventID       uint64
	ResultFilters []*trader.ResultFilter
	StatFilters   []*trader.StatFilter
}
