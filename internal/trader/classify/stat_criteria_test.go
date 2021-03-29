package classify

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseStatValue(t *testing.T) {
	t.Run("parses value from TeamStats struct", func(t *testing.T) {
		t.Helper()

		tc := []struct {
			Stats *statistico.TeamStats
			Stat  string
			Value uint32
		}{
			{
				Stats: &statistico.TeamStats{Goals: &wrappers.UInt32Value{Value: 4}},
				Stat:  "GOALS",
				Value: 4,
			},
			{
				Stats: &statistico.TeamStats{ShotsOnGoal: &wrappers.UInt32Value{Value: 12}},
				Stat:  "SHOTS_ON_GOAL",
				Value: 12,
			},
		}

		for _, c := range tc {
			val, err := parseStatValue(c.Stats, c.Stat)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.Equal(t, c.Value, val)
		}
	})

	t.Run("returns an error if stat provided is not supported", func(t *testing.T) {
		t.Helper()

		_, err := parseStatValue(&statistico.TeamStats{}, "INVALID")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "stat INVALID is not supported", err.Error())
	})
}

func Test_parseTeamStats(t *testing.T) {
	t.Run("returns a TeamStats struct associated to a team and action", func(t *testing.T) {
		t.Helper()

		tc := []struct {
			Result     *statistico.Result
			TeamID     uint64
			Action     string
			ExpectedID uint64
		}{
			{
				Result: &statistico.Result{
					HomeTeam:      &statistico.Team{Id: 55},
					HomeTeamStats: &statistico.TeamStats{TeamId: 55},
				},
				TeamID:     55,
				Action:     "FOR",
				ExpectedID: 55,
			},
			{
				Result: &statistico.Result{
					HomeTeam:      &statistico.Team{Id: 55},
					AwayTeamStats: &statistico.TeamStats{TeamId: 2},
				},
				TeamID:     55,
				Action:     "AGAINST",
				ExpectedID: 2,
			},
		}

		for _, c := range tc {
			st, err := parseTeamStats(c.Result, c.TeamID, c.Action)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.Equal(t, st.TeamId, c.ExpectedID)
		}
	})

	t.Run("returns an error if unable to parse team stats from a result", func(t *testing.T) {
		t.Helper()

		res := &statistico.Result{
			Id:            19281,
			HomeTeam:      &statistico.Team{Id: 55},
			HomeTeamStats: &statistico.TeamStats{TeamId: 55},
		}

		_, err := parseTeamStats(res, 55, "AGAINST")

		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		assert.Equal(t, "no stats available for team 55 and result 19281", err.Error())
	})

	t.Run("returns an error if action provided is not supported", func(t *testing.T) {
		t.Helper()

		_, err := parseTeamStats(&statistico.Result{}, 55, "INVALID")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "action INVALID is not supported", err.Error())
	})
}

func Test_parseStatValues(t *testing.T) {
	t.Run("returns a slice of uint32 values", func(t *testing.T) {
		t.Helper()

		tc := []struct {
			Results  []*statistico.Result
			TeamID   uint64
			Filter   *trader.StatFilter
			Expected []uint32
		}{
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 0}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:   "GOALS",
					Action: "FOR",
				},
				Expected: []uint32{4, 0, 1},
			},
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 0}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:   "GOALS",
					Action: "AGAINST",
				},
				Expected: []uint32{4, 0, 1},
			},
		}

		for _, c := range tc {
			values, err := parseStatValues(c.Results, c.TeamID, c.Filter)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.Equal(t, c.Expected, values)
		}
	})

	t.Run("returns an error if error returned by supporting method", func(t *testing.T) {
		t.Helper()

		results := []*statistico.Result{
			{
				Id:            1,
				HomeTeam:      &statistico.Team{Id: 55},
				HomeTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 4}},
			},
			{
				Id:            2,
				HomeTeam:      &statistico.Team{Id: 55},
				HomeTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 0}},
			},
			{
				Id:            3,
				HomeTeam:      &statistico.Team{Id: 55},
				AwayTeamStats: &statistico.TeamStats{TeamId: 2, Goals: &wrappers.UInt32Value{Value: 1}},
			},
		}

		filter := &trader.StatFilter{
			Stat:   "GOALS",
			Action: "FOR",
		}

		_, err := parseStatValues(results, 55, filter)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "no stats available for team 55 and result 3", err.Error())
	})
}

func Test_statMeetsCriteria(t *testing.T) {
	t.Run("returns true if results match provided filter", func(t *testing.T) {
		t.Helper()

		tc := []struct {
			Results []*statistico.Result
			TeamID  uint64
			Filter  *trader.StatFilter
		}{
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 0}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "GOALS",
					Action:  "FOR",
					Measure: "TOTAL",
					Metric:  "GTE",
					Value:   4,
				},
			},
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 1, Goals: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 2, Goals: &wrappers.UInt32Value{Value: 0}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 3, Goals: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "GOALS",
					Action:  "AGAINST",
					Measure: "TOTAL",
					Metric:  "GTE",
					Value:   5,
				},
			},
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 1, ShotsOnGoal: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 2, ShotsOnGoal: &wrappers.UInt32Value{Value: 0}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 3, ShotsOnGoal: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "SHOTS_ON_GOAL",
					Action:  "AGAINST",
					Measure: "AVERAGE",
					Metric:  "GTE",
					Value:   1.55,
				},
			},
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 1, ShotsOnGoal: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 2, ShotsOnGoal: &wrappers.UInt32Value{Value: 0}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 3, ShotsOnGoal: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "SHOTS_ON_GOAL",
					Action:  "AGAINST",
					Measure: "AVERAGE",
					Metric:  "LTE",
					Value:   2.15,
				},
			},
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 1, ShotsOnGoal: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 2, ShotsOnGoal: &wrappers.UInt32Value{Value: 5}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 3, ShotsOnGoal: &wrappers.UInt32Value{Value: 2}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "SHOTS_ON_GOAL",
					Action:  "AGAINST",
					Measure: "CONTINUOUS",
					Metric:  "GTE",
					Value:   1,
				},
			},
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 1, ShotsOnGoal: &wrappers.UInt32Value{Value: 1}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 2, ShotsOnGoal: &wrappers.UInt32Value{Value: 2}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 3, ShotsOnGoal: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "SHOTS_ON_GOAL",
					Action:  "AGAINST",
					Measure: "CONTINUOUS",
					Metric:  "LTE",
					Value:   2,
				},
			},
		}

		for _, c := range tc {
			yes, err := statMeetsCriteria(c.Results, c.TeamID, c.Filter)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.True(t, yes)
		}
	})

	t.Run("returns false if results do not match provided filter", func(t *testing.T) {
		t.Helper()

		tc := []struct {
			Results []*statistico.Result
			TeamID  uint64
			Filter  *trader.StatFilter
		}{
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 0}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 55, Goals: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "GOALS",
					Action:  "FOR",
					Measure: "TOTAL",
					Metric:  "GTE",
					Value:   6,
				},
			},
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 1, Goals: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 2, Goals: &wrappers.UInt32Value{Value: 0}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 3, Goals: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "GOALS",
					Action:  "AGAINST",
					Measure: "TOTAL",
					Metric:  "GTE",
					Value:   10,
				},
			},
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 1},
						AwayTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 1, ShotsOnGoal: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 2},
						AwayTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 2, ShotsOnGoal: &wrappers.UInt32Value{Value: 0}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 3},
						AwayTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 3, ShotsOnGoal: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "SHOTS_ON_GOAL",
					Action:  "AGAINST",
					Measure: "AVERAGE",
					Metric:  "GTE",
					Value:   2.5,
				},
			},
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 1, ShotsOnGoal: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 2, ShotsOnGoal: &wrappers.UInt32Value{Value: 0}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 3, ShotsOnGoal: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "SHOTS_ON_GOAL",
					Action:  "AGAINST",
					Measure: "AVERAGE",
					Metric:  "LTE",
					Value:   1.10,
				},
			},
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 1},
						AwayTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 1, ShotsOnGoal: &wrappers.UInt32Value{Value: 4}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 2},
						AwayTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 2, ShotsOnGoal: &wrappers.UInt32Value{Value: 0}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 3},
						AwayTeam:      &statistico.Team{Id: 55},
						HomeTeamStats: &statistico.TeamStats{TeamId: 3, ShotsOnGoal: &wrappers.UInt32Value{Value: 4}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "SHOTS_ON_GOAL",
					Action:  "AGAINST",
					Measure: "CONTINUOUS",
					Metric:  "GTE",
					Value:   2.5,
				},
			},
			{
				Results: []*statistico.Result{
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 1, ShotsOnGoal: &wrappers.UInt32Value{Value: 1}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 2, ShotsOnGoal: &wrappers.UInt32Value{Value: 2}},
					},
					{
						HomeTeam:      &statistico.Team{Id: 55},
						AwayTeamStats: &statistico.TeamStats{TeamId: 3, ShotsOnGoal: &wrappers.UInt32Value{Value: 1}},
					},
				},
				TeamID: 55,
				Filter: &trader.StatFilter{
					Stat:    "SHOTS_ON_GOAL",
					Action:  "AGAINST",
					Measure: "CONTINUOUS",
					Metric:  "LTE",
					Value:   1,
				},
			},
		}

		for _, c := range tc {
			yes, err := statMeetsCriteria(c.Results, c.TeamID, c.Filter)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.False(t, yes)
		}
	})
}
