package strategy_test

import (
	"context"
	"errors"
	"github.com/statistico/statistico-proto/go"
	m "github.com/statistico/statistico-trader/internal/trader/mock"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestStatClassifier_MatchesFilter(t *testing.T) {
	t.Run("returns true if result matches stat filter", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			Fixture        *strategy.Fixture
			FetchedResults []*statistico.Result
			Filter         *strategy.StatFilter
		}{
			{
				Fixture: &strategy.Fixture{
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
				Filter: &strategy.StatFilter{
					Stat:    "GOALS",
					Team:    "HOME_TEAM",
					Action:  "FOR",
					Games:   3,
					Measure: "TOTAL",
					Metric:  "GTE",
					Value:   3,
					Venue:   "HOME_AWAY",
				},
			},
			{
				Fixture: &strategy.Fixture{
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
				Filter: &strategy.StatFilter{
					Stat:    "GOALS",
					Team:    "HOME_TEAM",
					Action:  "FOR",
					Games:   3,
					Measure: "AVERAGE",
					Metric:  "LTE",
					Value:   3.5,
					Venue:   "HOME_AWAY",
				},
			},
		}

		for index, res := range assertions {
			client := new(m.ResultClient)
			classifier := strategy.NewStatFilterClassifier(client)

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

	t.Run("returns false if result matches stat filter", func(t *testing.T) {
		t.Helper()

		assertions := []struct {
			Fixture        *strategy.Fixture
			FetchedResults []*statistico.Result
			Filter         *strategy.StatFilter
		}{
			{
				Fixture: &strategy.Fixture{
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
				Filter: &strategy.StatFilter{
					Stat:    "GOALS",
					Team:    "HOME_TEAM",
					Action:  "FOR",
					Games:   3,
					Measure: "TOTAL",
					Metric:  "GTE",
					Value:   15,
					Venue:   "HOME_AWAY",
				},
			},
			{
				Fixture: &strategy.Fixture{
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
				Filter: &strategy.StatFilter{
					Stat:    "GOALS",
					Team:    "HOME_TEAM",
					Action:  "FOR",
					Games:   3,
					Measure: "AVERAGE",
					Metric:  "LTE",
					Value:   2,
					Venue:   "HOME_AWAY",
				},
			},
		}

		for index, res := range assertions {
			client := new(m.ResultClient)
			classifier := strategy.NewStatFilterClassifier(client)

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
		classifier := strategy.NewStatFilterClassifier(client)

		fixture := &strategy.Fixture{
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

		filter := &strategy.StatFilter{
			Stat:    "GOALS",
			Team:    "AWAY_TEAM",
			Action:  "FOR",
			Games:   3,
			Measure: "AVERAGE",
			Metric:  "LTE",
			Value:   2,
			Venue:   "HOME_AWAY",
		}

		req := mock.MatchedBy(func(r *statistico.TeamResultRequest) bool {
			a := assert.New(t)
			a.Equal(uint64(2), r.TeamId)
			a.Equal([]uint64{8}, r.SeasonIds)
			a.Equal(uint64(filter.Games), r.GetLimit().GetValue())
			a.Equal(fixture.Date.Format(time.RFC3339), r.GetDateBefore().GetValue())
			a.Equal("HOME_AWAY", r.GetVenue().GetValue())
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
