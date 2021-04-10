package classify_test

import (
	"context"
	"errors"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader/classify"
	mock2 "github.com/statistico/statistico-trader/internal/trader/mock"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestFilterMatcher_MatchesFilters(t *testing.T) {
	f1 := &strategy.ResultFilter{
		Team:   "HOME",
		Result: "WIN",
		Games:  3,
		Venue:  "HOME",
	}

	f2 := &strategy.ResultFilter{
		Team:   "AWAY",
		Result: "LOSE_DRAW",
		Games:  5,
		Venue:  "AWAY",
	}

	f3 := &strategy.StatFilter{
		Stat:    "GOALS",
		Team:    "HOME_TEAM",
		Action:  "FOR",
		Games:   3,
		Measure: "AVG",
		Metric:  "LTE",
		Value:   3.5,
		Venue:   "HOME_AWAY",
	}

	f4 := &strategy.StatFilter{
		Stat:    "GOALS",
		Team:    "HOME_TEAM",
		Action:  "FOR",
		Games:   3,
		Measure: "AVG",
		Metric:  "LTE",
		Value:   2,
		Venue:   "HOME_AWAY",
	}

	ctx := context.Background()

	results := []*strategy.ResultFilter{f1, f2}
	stats := []*strategy.StatFilter{f3, f4}

	query := classify.MatcherQuery{
		EventID:       192810,
		ResultFilters: results,
		StatFilters:   stats,
	}

	fixture := statistico.Fixture{
		Id:          192810,
		Competition: &statistico.Competition{Id: 8},
		Season:      &statistico.Season{Id: 17420},
		HomeTeam:    &statistico.Team{Id: 5},
		AwayTeam:    &statistico.Team{Id: 10},
		DateTime:    &statistico.Date{Utc: 1616052304},
	}

	fix := classify.Fixture{
		ID:         192810,
		HomeTeamID: 5,
		AwayTeamID: 10,
		Date:       time.Unix(1616052304, 0),
		SeasonID:   17420,
	}

	t.Run("returns bool if Fixture matches all filters provided", func(t *testing.T) {
		t.Helper()

		fc := new(mock2.FixtureClient)
		rc := new(MockResultClassifier)
		sc := new(MockStatClassifier)

		fc.On("ByID", ctx, uint64(192810)).Return(&fixture, nil)

		rc.On("MatchesFilter", ctx, &fix, f1).Return(true, nil)
		rc.On("MatchesFilter", ctx, &fix, f2).Return(true, nil)

		sc.On("MatchesFilter", ctx, &fix, f3).Return(true, nil)
		sc.On("MatchesFilter", ctx, &fix, f4).Return(true, nil)

		matcher := classify.NewFilterMatcher(fc, rc, sc)

		matches, err := matcher.MatchesFilters(ctx, &query)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.True(t, matches)
		rc.AssertExpectations(t)
		sc.AssertExpectations(t)
	})

	t.Run("returns false if result classifier returns false for a filter", func(t *testing.T) {
		t.Helper()

		fc := new(mock2.FixtureClient)
		rc := new(MockResultClassifier)
		sc := new(MockStatClassifier)

		fc.On("ByID", ctx, uint64(192810)).Return(&fixture, nil)

		rc.On("MatchesFilter", ctx, &fix, f1).Return(true, nil)
		rc.On("MatchesFilter", ctx, &fix, f2).Return(false, nil)

		sc.AssertNotCalled(t, "MatchesFilter")

		matcher := classify.NewFilterMatcher(fc, rc, sc)

		matches, err := matcher.MatchesFilters(ctx, &query)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.False(t, matches)
		rc.AssertExpectations(t)
		sc.AssertExpectations(t)
	})

	t.Run("returns false if stat classifier returns false for a filter", func(t *testing.T) {
		t.Helper()

		fc := new(mock2.FixtureClient)
		rc := new(MockResultClassifier)
		sc := new(MockStatClassifier)

		fc.On("ByID", ctx, uint64(192810)).Return(&fixture, nil)

		rc.On("MatchesFilter", ctx, &fix, f1).Return(true, nil)
		rc.On("MatchesFilter", ctx, &fix, f2).Return(true, nil)

		sc.On("MatchesFilter", ctx, &fix, f3).Return(true, nil)
		sc.On("MatchesFilter", ctx, &fix, f4).Return(false, nil)

		matcher := classify.NewFilterMatcher(fc, rc, sc)

		matches, err := matcher.MatchesFilters(ctx, &query)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.False(t, matches)
		rc.AssertExpectations(t)
		sc.AssertExpectations(t)
	})

	t.Run("returns error if returned by result classifier", func(t *testing.T) {
		t.Helper()

		fc := new(mock2.FixtureClient)
		rc := new(MockResultClassifier)
		sc := new(MockStatClassifier)

		fc.On("ByID", ctx, uint64(192810)).Return(&fixture, nil)

		e := errors.New("error from classifier")

		rc.On("MatchesFilter", ctx, &fix, f1).Return(true, nil)
		rc.On("MatchesFilter", ctx, &fix, f2).Return(false, e)

		sc.AssertNotCalled(t, "MatchesFilter")

		matcher := classify.NewFilterMatcher(fc, rc, sc)

		_, err := matcher.MatchesFilters(ctx, &query)

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, e, err)
		rc.AssertExpectations(t)
		sc.AssertExpectations(t)
	})

	t.Run("returns error if returned by stat classifier", func(t *testing.T) {
		t.Helper()

		fc := new(mock2.FixtureClient)
		rc := new(MockResultClassifier)
		sc := new(MockStatClassifier)

		fc.On("ByID", ctx, uint64(192810)).Return(&fixture, nil)

		e := errors.New("error from classifier")

		rc.On("MatchesFilter", ctx, &fix, f1).Return(true, nil)
		rc.On("MatchesFilter", ctx, &fix, f2).Return(true, nil)

		sc.On("MatchesFilter", ctx, &fix, f3).Return(true, nil)
		sc.On("MatchesFilter", ctx, &fix, f4).Return(false, e)

		matcher := classify.NewFilterMatcher(fc, rc, sc)

		_, err := matcher.MatchesFilters(ctx, &query)

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, e, err)
		rc.AssertExpectations(t)
		sc.AssertExpectations(t)
	})
}

type MockResultClassifier struct {
	mock.Mock
}

func (m *MockResultClassifier) MatchesFilter(ctx context.Context, fix *classify.Fixture, f *strategy.ResultFilter) (bool, error) {
	args := m.Called(ctx, fix, f)
	return args.Get(0).(bool), args.Error(1)
}

type MockStatClassifier struct {
	mock.Mock
}

func (m *MockStatClassifier) MatchesFilter(ctx context.Context, fix *classify.Fixture, f *strategy.StatFilter) (bool, error) {
	args := m.Called(ctx, fix, f)
	return args.Get(0).(bool), args.Error(1)
}
