package classify_test

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/statistico/statistico-strategy/internal/trader/classify"
	m "github.com/statistico/statistico-strategy/internal/trader/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestResultClassifier_MatchesFilter(t *testing.T) {
	t.Run("returns true if result matches result filter", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			Fixture        *classify.Fixture
			FetchedResults []*statistico.Result
			Filter         *trader.ResultFilter
		}{
			{
				Fixture: &classify.Fixture{
					ID:         55,
					HomeTeamID: 1,
					AwayTeamID: 2,
					Date:       time.Unix(1584014400, 0),
					SeasonID:   8,
				},
				FetchedResults: []*statistico.Result{
					newProtoResult(1, 50, 4, 1),
					newProtoResult(10, 1, 1, 4),
					newProtoResult(1, 10, 2, 0),
				},
				Filter: &trader.ResultFilter{
					Team:   "HOME_TEAM",
					Result: "WIN",
					Games:  3,
					Venue:  "HOME",
				},
			},
			{
				Fixture: &classify.Fixture{
					ID:         55,
					HomeTeamID: 1,
					AwayTeamID: 2,
					Date:       time.Unix(1584014400, 0),
					SeasonID:   8,
				},
				FetchedResults: []*statistico.Result{
					newProtoResult(1, 5, 4, 4),
					newProtoResult(10, 1, 1, 4),
					newProtoResult(1, 11, 2, 2),
				},
				Filter: &trader.ResultFilter{
					Team:   "HOME_TEAM",
					Result: "WIN_DRAW",
					Games:  3,
					Venue:  "HOME",
				},
			},
			{
				Fixture: &classify.Fixture{
					ID:         55,
					HomeTeamID: 1,
					AwayTeamID: 2,
					Date:       time.Unix(1584014400, 0),
					SeasonID:   8,
				},
				FetchedResults: []*statistico.Result{
					newProtoResult(1, 5, 2, 4),
					newProtoResult(10, 1, 5, 4),
					newProtoResult(1, 111, 2, 4),
				},
				Filter: &trader.ResultFilter{
					Team:   "HOME_TEAM",
					Result: "LOSE",
					Games:  3,
					Venue:  "HOME",
				},
			},
			{
				Fixture: &classify.Fixture{
					ID:         55,
					HomeTeamID: 1,
					AwayTeamID: 2,
					Date:       time.Unix(1584014400, 0),
					SeasonID:   8,
				},
				FetchedResults: []*statistico.Result{
					newProtoResult(1, 5, 4, 4),
					newProtoResult(10, 1, 5, 4),
					newProtoResult(1, 11, 2, 4),
				},
				Filter: &trader.ResultFilter{
					Team:   "HOME_TEAM",
					Result: "LOSE_DRAW",
					Games:  3,
					Venue:  "HOME",
				},
			},
			{
				Fixture: &classify.Fixture{
					ID:         55,
					HomeTeamID: 1,
					AwayTeamID: 2,
					Date:       time.Unix(1584014400, 0),
					SeasonID:   8,
				},
				FetchedResults: []*statistico.Result{
					newProtoResult(1, 5, 4, 4),
					newProtoResult(10, 1, 5, 5),
					newProtoResult(1, 11, 2, 2),
				},
				Filter: &trader.ResultFilter{
					Team:   "HOME_TEAM",
					Result: "DRAW",
					Games:  3,
					Venue:  "HOME",
				},
			},
			{
				Fixture: &classify.Fixture{
					ID:         55,
					HomeTeamID: 5,
					AwayTeamID: 1,
					Date:       time.Unix(1584014400, 0),
					SeasonID:   8,
				},
				FetchedResults: []*statistico.Result{
					newProtoResult(1, 5, 4, 4),
					newProtoResult(5, 1, 5, 1),
					newProtoResult(1, 11, 2, 2),
				},
				Filter: &trader.ResultFilter{
					Team:   "AWAY_TEAM",
					Result: "LOSE_DRAW",
					Games:  3,
					Venue:  "HOME_AWAY",
				},
			},
		}

		for index, res := range assertions {
			client := new(m.ResultClient)
			classifier := classify.NewResultFilterClassifier(client)

			ctx := context.Background()

			req := mock.MatchedBy(func(r *statistico.TeamResultRequest) bool {
				a := assert.New(t)
				a.Equal(uint64(1), r.TeamId)
				a.Equal([]uint64{8}, r.SeasonIds)
				a.Equal(uint64(res.Filter.Games), r.GetLimit().GetValue())
				a.Equal(res.Fixture.Date.Format(time.RFC3339), r.GetDateBefore().GetValue())
				a.Equal(res.Filter.Venue, r.GetVenue().GetValue())
				return true
			})

			client.On("ByTeam", ctx, req).Return(res.FetchedResults, nil)

			success, err := classifier.MatchesFilter(ctx, res.Fixture, res.Filter)

			if err != nil {
				t.Fatalf("Expected nil, got %s at index %d", err.Error(), index)
			}

			assert.True(t, success)
		}
	})

	t.Run("returns false if result does not match filter", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			Fixture        *classify.Fixture
			FetchedResults []*statistico.Result
			Filter         *trader.ResultFilter
		}{
			{
				Fixture: &classify.Fixture{
					ID:         55,
					HomeTeamID: 1,
					AwayTeamID: 2,
					Date:       time.Unix(1584014400, 0),
					SeasonID:   8,
				},
				FetchedResults: []*statistico.Result{
					newProtoResult(1, 50, 4, 1),
					newProtoResult(10, 1, 1, 4),
					newProtoResult(1, 10, 2, 2),
				},
				Filter: &trader.ResultFilter{
					Team:   "HOME_TEAM",
					Result: "WIN",
					Games:  3,
					Venue:  "HOME",
				},
			},
			{
				Fixture: &classify.Fixture{
					ID:         55,
					HomeTeamID: 1,
					AwayTeamID: 2,
					Date:       time.Unix(1584014400, 0),
					SeasonID:   8,
				},
				FetchedResults: []*statistico.Result{
					newProtoResult(1, 5, 4, 4),
					newProtoResult(10, 1, 5, 4),
					newProtoResult(1, 11, 2, 2),
				},
				Filter: &trader.ResultFilter{
					Team:   "HOME_TEAM",
					Result: "WIN_DRAW",
					Games:  3,
					Venue:  "HOME",
				},
			},
			{
				Fixture: &classify.Fixture{
					ID:         55,
					HomeTeamID: 1,
					AwayTeamID: 2,
					Date:       time.Unix(1584014400, 0),
					SeasonID:   8,
				},
				FetchedResults: []*statistico.Result{
					newProtoResult(1, 5, 4, 4),
					newProtoResult(10, 1, 5, 4),
					newProtoResult(1, 111, 2, 4),
				},
				Filter: &trader.ResultFilter{
					Team:   "HOME_TEAM",
					Result: "LOSE",
					Games:  3,
					Venue:  "HOME",
				},
			},
			{
				Fixture: &classify.Fixture{
					ID:         55,
					HomeTeamID: 1,
					AwayTeamID: 2,
					Date:       time.Unix(1584014400, 0),
					SeasonID:   8,
				},
				FetchedResults: []*statistico.Result{
					newProtoResult(1, 5, 4, 4),
					newProtoResult(10, 1, 5, 4),
					newProtoResult(1, 11, 5, 4),
				},
				Filter: &trader.ResultFilter{
					Team:   "HOME_TEAM",
					Result: "LOSE_DRAW",
					Games:  3,
					Venue:  "HOME",
				},
			},
			{
				Fixture: &classify.Fixture{
					ID:         55,
					HomeTeamID: 1,
					AwayTeamID: 2,
					Date:       time.Unix(1584014400, 0),
					SeasonID:   8,
				},
				FetchedResults: []*statistico.Result{
					newProtoResult(1, 5, 2, 4),
					newProtoResult(10, 1, 5, 5),
					newProtoResult(1, 11, 2, 2),
				},
				Filter: &trader.ResultFilter{
					Team:   "HOME_TEAM",
					Result: "DRAW",
					Games:  3,
					Venue:  "HOME",
				},
			},
		}

		for index, res := range assertions {
			client := new(m.ResultClient)
			classifier := classify.NewResultFilterClassifier(client)

			ctx := context.Background()

			req := mock.MatchedBy(func(r *statistico.TeamResultRequest) bool {
				a := assert.New(t)
				a.Equal(uint64(1), r.TeamId)
				a.Equal([]uint64{8}, r.SeasonIds)
				a.Equal(uint64(res.Filter.Games), r.GetLimit().GetValue())
				a.Equal(res.Fixture.Date.Format(time.RFC3339), r.GetDateBefore().GetValue())
				a.Equal(res.Filter.Venue, r.GetVenue().GetValue())
				return true
			})

			client.On("ByTeam", ctx, req).Return(res.FetchedResults, nil)

			success, err := classifier.MatchesFilter(ctx, res.Fixture, res.Filter)

			if err != nil {
				t.Fatalf("Expected nil, got %s at index %d", err.Error(), index)
			}

			assert.False(t, success)
		}
	})

	t.Run("returns false and an error if error is sent in error channel", func(t *testing.T) {
		t.Helper()

		client := new(m.ResultClient)
		classifier := classify.NewResultFilterClassifier(client)

		fixture := &classify.Fixture{
			ID:         55,
			HomeTeamID: 1,
			AwayTeamID: 2,
			Date:       time.Unix(1584014400, 0),
			SeasonID:   8,
		}

		results := []*statistico.Result{
			newProtoResult(1, 5, 4, 4),
			newProtoResult(10, 1, 5, 5),
			newProtoResult(1, 11, 2, 2),
		}

		ctx := context.Background()

		filter := &trader.ResultFilter{
			Team:   "AWAY_TEAM",
			Result: "DRAW",
			Games:  3,
			Venue:  "HOME",
		}

		req := mock.MatchedBy(func(r *statistico.TeamResultRequest) bool {
			a := assert.New(t)
			a.Equal(uint64(2), r.TeamId)
			a.Equal([]uint64{8}, r.SeasonIds)
			a.Equal(uint64(filter.Games), r.GetLimit().GetValue())
			a.Equal(fixture.Date.Format(time.RFC3339), r.GetDateBefore().GetValue())
			a.Equal("HOME", r.GetVenue().GetValue())
			return true
		})

		client.On("ByTeam", ctx, req).Return(results, errors.New("invalid argument"))

		success, err := classifier.MatchesFilter(ctx, fixture, filter)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.False(t, success)
		assert.Equal(t, "invalid argument", err.Error())
	})
}

func newProtoResult(homeId, awayId uint64, homeScore, awayScore uint32) *statistico.Result {
	return &statistico.Result{
		Id:       1,
		HomeTeam: &statistico.Team{Id: homeId},
		AwayTeam: &statistico.Team{Id: awayId},
		Stats: &statistico.MatchStats{
			HomeScore: &wrappers.UInt32Value{Value: homeScore},
			AwayScore: &wrappers.UInt32Value{Value: awayScore},
		},
		HomeTeamStats: &statistico.TeamStats{Goals: &wrappers.UInt32Value{Value: homeScore}},
		AwayTeamStats: &statistico.TeamStats{Goals: &wrappers.UInt32Value{Value: awayScore}},
	}
}
