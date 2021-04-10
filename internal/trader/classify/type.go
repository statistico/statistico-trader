package classify

import (
	"github.com/statistico/statistico-trader/internal/trader/strategy"
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
	ResultFilters []*strategy.ResultFilter
	StatFilters   []*strategy.StatFilter
}
