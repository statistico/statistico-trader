package grpc

import (
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader/market"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_transformTradeResultToStrategyTrade(t *testing.T) {
	t.Run("transforms trade Result struct into StrategyTrade struct", func(t *testing.T) {
		t.Helper()

		result := "SUCCESS"

		tr := market.Trade{
			MarketName:    "MATCH_ODDS",
			RunnerName:    "Draw",
			RunnerPrice:   4.05,
			EventId:       198187871,
			CompetitionId: 8,
			SeasonId:      17420,
			EventDate:     time.Unix(1584014400, 0),
			Side:          "BACK",
			Result:        &result,
		}

		st, err := transformTradeResultToStrategyTrade(&tr)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		a := assert.New(t)

		a.Equal("MATCH_ODDS", st.MarketName)
		a.Equal("Draw", st.RunnerName)
		a.Equal(float32(4.05), st.RunnerPrice)
		a.Equal(statistico.SideEnum_BACK, st.Side)
		a.Equal(uint64(198187871), st.EventId)
		a.Equal(uint64(8), st.CompetitionId)
		a.Equal(uint64(17420), st.SeasonId)
		a.Equal(int64(1584014400), st.EventDate.Seconds)
		a.Equal(statistico.TradeResultEnum_SUCCESS, st.Result)
	})
}

func Test_transformResultFilters(t *testing.T) {
	t.Run("transforms statistico.ResultFilter struct into trade.ResultFilter struct", func(t *testing.T) {
		t.Helper()

		f := []*statistico.ResultFilter{
			{
				Team:   statistico.TeamEnum_HOME_TEAM,
				Result: statistico.ResultEnum_LOSE_DRAW,
				Games:  4,
				Venue:  statistico.VenueEnum_AWAY,
			},
			{
				Team:   statistico.TeamEnum_AWAY_TEAM,
				Result: statistico.ResultEnum_LOSE,
				Games:  6,
				Venue:  statistico.VenueEnum_AWAY,
			},
		}

		fs := transformResultFilters(f)

		a := assert.New(t)
		a.Equal(2, len(fs))
		a.Equal("HOME_TEAM", fs[0].Team)
		a.Equal("LOSE_DRAW", fs[0].Result)
		a.Equal(uint8(4), fs[0].Games)
		a.Equal("AWAY", fs[0].Venue)
		a.Equal("AWAY_TEAM", fs[1].Team)
		a.Equal("LOSE", fs[1].Result)
		a.Equal(uint8(6), fs[1].Games)
		a.Equal("AWAY", fs[1].Venue)
	})
}

func Test_transformStatFilters(t *testing.T) {
	t.Run("transforms statistico.StatFilter struct into trade.StatFilter struct", func(t *testing.T) {
		t.Helper()

		f := []*statistico.StatFilter{
			{
				Stat:    statistico.StatEnum_GOALS,
				Team:    statistico.TeamEnum_HOME_TEAM,
				Action:  statistico.ActionEnum_FOR,
				Games:   3,
				Measure: statistico.MeasureEnum_TOTAL,
				Metric:  statistico.MetricEnum_GTE,
				Value:   2.3,
				Venue:   statistico.VenueEnum_HOME_AWAY,
			},
			{
				Stat:    statistico.StatEnum_SHOTS_ON_GOAL,
				Team:    statistico.TeamEnum_HOME_TEAM,
				Action:  statistico.ActionEnum_AGAINST,
				Games:   5,
				Measure: statistico.MeasureEnum_TOTAL,
				Metric:  statistico.MetricEnum_LTE,
				Value:   10,
				Venue:   statistico.VenueEnum_HOME,
			},
		}

		fs := transformStatFilters(f)

		a := assert.New(t)
		a.Equal(2, len(fs))
		a.Equal("GOALS", fs[0].Stat)
		a.Equal("HOME_TEAM", fs[0].Team)
		a.Equal("FOR", fs[0].Action)
		a.Equal(uint8(3), fs[0].Games)
		a.Equal("TOTAL", fs[0].Measure)
		a.Equal("GTE", fs[0].Metric)
		a.Equal(float32(2.3), fs[0].Value)
		a.Equal("HOME_AWAY", fs[0].Venue)
		a.Equal("SHOTS_ON_GOAL", fs[1].Stat)
		a.Equal("HOME_TEAM", fs[1].Team)
		a.Equal("AGAINST", fs[1].Action)
		a.Equal(uint8(5), fs[1].Games)
		a.Equal("TOTAL", fs[1].Measure)
		a.Equal("LTE", fs[1].Metric)
		a.Equal(float32(10), fs[1].Value)
		a.Equal("HOME", fs[1].Venue)
	})
}
