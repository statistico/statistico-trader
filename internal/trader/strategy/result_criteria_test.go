package strategy

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-proto/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_meetsWinCriteria(t *testing.T) {
	t.Run("return true if result meets win criteria", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			TeamID uint64
			Result *statistico.Result
		}{
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
			{
				TeamID: 2,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 4}},
				},
			},
		}

		for _, res := range assertions {
			assert.True(t, meetsWinCriteria(res.TeamID, res.Result))
		}
	})

	t.Run("return false if result does not meet win criteria", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			TeamID uint64
			Result *statistico.Result
		}{
			{
				TeamID: 2,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 4}},
				},
			},
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 0},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
		}

		for _, res := range assertions {
			assert.False(t, meetsWinCriteria(res.TeamID, res.Result))
		}
	})
}

func Test_meetsWinDrawCriteria(t *testing.T) {
	t.Run("return true if result meets win draw criteria", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			TeamID uint64
			Result *statistico.Result
		}{
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 2}},
				},
			},
			{
				TeamID: 2,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 4}},
				},
			},
		}

		for _, res := range assertions {
			assert.True(t, meetsWinDrawCriteria(res.TeamID, res.Result))
		}
	})

	t.Run("return false if does not meet win draw criteria", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			TeamID uint64
			Result *statistico.Result
		}{
			{
				TeamID: 2,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 4}},
				},
			},
		}

		for _, res := range assertions {
			assert.False(t, meetsWinCriteria(res.TeamID, res.Result))
		}
	})
}

func Test_meetsLoseCriteria(t *testing.T) {
	t.Run("return true if result meets lose criteria", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			TeamID uint64
			Result *statistico.Result
		}{
			{
				TeamID: 2,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 4}},
				},
			},
		}

		for _, res := range assertions {
			assert.True(t, meetsLoseCriteria(res.TeamID, res.Result))
		}
	})

	t.Run("return false if result does not meet lose criteria", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			TeamID uint64
			Result *statistico.Result
		}{
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
			{
				TeamID: 2,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 4}},
				},
			},
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 0},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
		}

		for _, res := range assertions {
			assert.False(t, meetsLoseCriteria(res.TeamID, res.Result))
		}
	})
}

func Test_meetsLoseDrawCriteria(t *testing.T) {
	t.Run("return true if result meets lose draw criteria", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			TeamID uint64
			Result *statistico.Result
		}{
			{
				TeamID: 2,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 2}},
				},
			},
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 4}},
				},
			},
		}

		for _, res := range assertions {
			assert.True(t, meetsLoseDrawCriteria(res.TeamID, res.Result))
		}
	})

	t.Run("return false if result does not meet lose draw criteria", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			TeamID uint64
			Result *statistico.Result
		}{
			{
				TeamID: 1,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
			{
				TeamID: 2,
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 4}},
				},
			},
		}

		for _, res := range assertions {
			assert.False(t, meetsLoseDrawCriteria(res.TeamID, res.Result))
		}
	})
}

func Test_meetsWinLoseCriteria(t *testing.T) {
	t.Run("return true if result meets win lose criteria", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			TeamID uint64
			Result *statistico.Result
		}{
			{
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
			{
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 3}},
				},
			},
			{
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 2},
						AwayScore: &wrappers.UInt32Value{Value: 4}},
				},
			},
		}

		for _, res := range assertions {
			assert.True(t, meetsWinLoseCriteria(res.Result))
		}
	})

	t.Run("return false if result does not meet win lose criteria", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			TeamID uint64
			Result *statistico.Result
		}{
			{
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 0},
						AwayScore: &wrappers.UInt32Value{Value: 0}},
				},
			},
			{
				Result: &statistico.Result{
					HomeTeam: &statistico.Team{Id: 1},
					AwayTeam: &statistico.Team{Id: 2},
					Stats: &statistico.MatchStats{
						HomeScore: &wrappers.UInt32Value{Value: 4},
						AwayScore: &wrappers.UInt32Value{Value: 4}},
				},
			},
		}

		for _, res := range assertions {
			assert.False(t, meetsWinLoseCriteria(res.Result))
		}
	})
}
