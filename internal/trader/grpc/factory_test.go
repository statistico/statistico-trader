package grpc

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStrategyFromRequest(t *testing.T) {
	t.Run("convert a statistico SaveStrategyRequest into a trader Strategy struct", func(t *testing.T) {
		t.Helper()

		r := &statistico.SaveStrategyRequest{
			Name:           "Money Maker v1",
			Description:    "Home favourite strategy",
			UserId:         "a5f04fd2-dfe7-41c1-af38-d490119705d8",
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &wrappers.FloatValue{Value: 1.50},
			MaxOdds:        &wrappers.FloatValue{Value: 5.25},
			Side:           statistico.SideEnum_BACK,
			CompetitionIds: []uint64{8, 14},
			ResultFilters: []*statistico.ResultFilter{
				{
					Team:   statistico.TeamEnum_HOME_TEAM,
					Result: statistico.ResultEnum_WIN_DRAW,
					Games:  2,
					Venue:  statistico.VenueEnum_HOME_AWAY,
				},
			},
			StatFilters: []*statistico.StatFilter{
				{
					Stat:    statistico.StatEnum_GOALS,
					Team:    statistico.TeamEnum_HOME_TEAM,
					Action:  statistico.ActionEnum_AGAINST,
					Games:   4,
					Metric:  statistico.MetricEnum_GTE,
					Measure: statistico.MeasureEnum_AVERAGE,
					Value:   3.1,
					Venue:   statistico.VenueEnum_AWAY,
				},
			},
			Visibility: statistico.VisibilityEnum_PRIVATE,
			StakingPlan: &statistico.StakingPlan{
				Name:  statistico.StakingPlanEnum_PERCENTAGE,
				Value: 2.5,
			},
		}

		s, err := strategyFromRequest(r, time.Unix(1616936636, 0))

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		res := []*trader.ResultFilter{
			{
				Team:   "HOME_TEAM",
				Result: "WIN_DRAW",
				Games:  2,
				Venue:  "HOME_AWAY",
			},
		}

		stat := []*trader.StatFilter{
			{
				Stat:    "GOALS",
				Team:    "HOME_TEAM",
				Action:  "AGAINST",
				Games:   4,
				Measure: "AVERAGE",
				Metric:  "GTE",
				Value:   3.1,
				Venue:   "AWAY",
			},
		}

		plan := trader.StakingPlan{
			Name:   "PERCENTAGE",
			Number: 2.5,
		}

		a := assert.New(t)

		a.Equal(r.GetName(), s.Name)
		a.Equal(r.GetDescription(), s.Description)
		a.Equal(r.GetMarket(), s.MarketName)
		a.Equal(r.GetRunner(), s.RunnerName)
		a.Equal(r.GetCompetitionIds(), s.CompetitionIDs)
		a.Equal(r.GetSide().String(), s.Side)
		a.Equal(r.GetVisibility().String(), s.Visibility)
		a.Equal("ACTIVE", s.Status)
		a.Equal(plan, s.StakingPlan)
		a.Equal(res, s.ResultFilters)
		a.Equal(stat, s.StatFilters)
		a.Equal(time.Unix(1616936636, 0), s.CreatedAt)
		a.Equal(time.Unix(1616936636, 0), s.UpdatedAt)
	})

	t.Run("returns error if User ID is not a valid uuid string", func(t *testing.T) {
		t.Helper()

		r := &statistico.SaveStrategyRequest{
			Name:           "Money Maker v1",
			Description:    "Home favourite strategy",
			UserId:         "a5f04fd2",
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &wrappers.FloatValue{Value: 1.50},
			MaxOdds:        &wrappers.FloatValue{Value: 5.25},
			Side:           statistico.SideEnum_BACK,
			CompetitionIds: []uint64{8, 14},
			Visibility:     statistico.VisibilityEnum_PRIVATE,
		}

		_, err := strategyFromRequest(r, time.Unix(1616936636, 0))

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = InvalidArgument desc = error parsing user ID: invalid UUID length: 8", err.Error())
	})

	t.Run("returns error if both min and max odds are not provided", func(t *testing.T) {
		t.Helper()

		r := &statistico.SaveStrategyRequest{
			Name:           "Money Maker v1",
			Description:    "Home favourite strategy",
			UserId:         "a5f04fd2-dfe7-41c1-af38-d490119705d8",
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			Side:           statistico.SideEnum_BACK,
			CompetitionIds: []uint64{8, 14},
			Visibility:     statistico.VisibilityEnum_PRIVATE,
			StakingPlan: &statistico.StakingPlan{
				Name:  statistico.StakingPlanEnum_PERCENTAGE,
				Value: 2.5,
			},
		}

		_, err := strategyFromRequest(r, time.Unix(1616936636, 0))

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = InvalidArgument desc = Min and max odds cannot both be nil", err.Error())
	})

	t.Run("returns error if staking plan value equal to or less than zero", func(t *testing.T) {
		t.Helper()

		r := &statistico.SaveStrategyRequest{
			Name:           "Money Maker v1",
			Description:    "Home favourite strategy",
			UserId:         "a5f04fd2-dfe7-41c1-af38-d490119705d8",
			Market:         "MATCH_ODDS",
			Runner:         "Home",
			MinOdds:        &wrappers.FloatValue{Value: 1.50},
			MaxOdds:        &wrappers.FloatValue{Value: 5.25},
			Side:           statistico.SideEnum_BACK,
			CompetitionIds: []uint64{8, 14},
			ResultFilters: []*statistico.ResultFilter{
				{
					Team:   statistico.TeamEnum_HOME_TEAM,
					Result: statistico.ResultEnum_WIN_DRAW,
					Games:  2,
					Venue:  statistico.VenueEnum_HOME_AWAY,
				},
			},
			StatFilters: []*statistico.StatFilter{
				{
					Stat:    statistico.StatEnum_GOALS,
					Team:    statistico.TeamEnum_HOME_TEAM,
					Action:  statistico.ActionEnum_AGAINST,
					Games:   4,
					Metric:  statistico.MetricEnum_GTE,
					Measure: statistico.MeasureEnum_AVERAGE,
					Value:   3.1,
					Venue:   statistico.VenueEnum_AWAY,
				},
			},
			Visibility: statistico.VisibilityEnum_PRIVATE,
			StakingPlan: &statistico.StakingPlan{
				Name:  statistico.StakingPlanEnum_PERCENTAGE,
				Value: 0,
			},
		}

		_, err := strategyFromRequest(r, time.Unix(1616936636, 0))

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = InvalidArgument desc = staking plan must be greater than zero", err.Error())
	})
}
