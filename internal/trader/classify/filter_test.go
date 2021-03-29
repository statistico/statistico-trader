package classify_test

import (
	"context"
	"errors"
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/statistico/statistico-strategy/internal/trader/classify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFilterMatcher_MatchesFilters(t *testing.T) {
	f1 := &trader.ResultFilter{
		Team:   "HOME",
		Result: "WIN",
		Games:  3,
		Venue:  "HOME",
	}

	f2 := &trader.ResultFilter{
		Team:   "AWAY",
		Result: "LOSE_DRAW",
		Games:  5,
		Venue:  "AWAY",
	}

	f3 := &trader.StatFilter{
		Stat:    "GOALS",
		Team:    "HOME_TEAM",
		Action:  "FOR",
		Games:   3,
		Measure: "AVG",
		Metric:  "LTE",
		Value:   3.5,
		Venue:   "HOME_AWAY",
	}

	f4 := &trader.StatFilter{
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

	fix := &classify.Fixture{}

	results := []*trader.ResultFilter{f1, f2}
	stats := []*trader.StatFilter{f3, f4}

	t.Run("returns bool if Fixture matches all filters provided", func(t *testing.T) {
		t.Helper()

		rc := new(MockResultClassifier)
		sc := new(MockStatClassifier)

		rc.On("MatchesFilter", ctx, fix, f1).Return(true, nil)
		rc.On("MatchesFilter", ctx, fix, f2).Return(true, nil)

		sc.On("MatchesFilter", ctx, fix, f3).Return(true, nil)
		sc.On("MatchesFilter", ctx, fix, f4).Return(true, nil)

		matcher := classify.NewFilterMatcher(rc, sc)

		matches, err := matcher.MatchesFilters(ctx, fix, results, stats)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.True(t, matches)
		rc.AssertExpectations(t)
		sc.AssertExpectations(t)
	})

	t.Run("returns false if result classifier returns false for a filter", func(t *testing.T) {
		t.Helper()

		rc := new(MockResultClassifier)
		sc := new(MockStatClassifier)

		rc.On("MatchesFilter", ctx, fix, f1).Return(true, nil)
		rc.On("MatchesFilter", ctx, fix, f2).Return(false, nil)

		sc.AssertNotCalled(t, "MatchesFilter")

		matcher := classify.NewFilterMatcher(rc, sc)

		matches, err := matcher.MatchesFilters(ctx, fix, results, stats)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.False(t, matches)
		rc.AssertExpectations(t)
		sc.AssertExpectations(t)
	})

	t.Run("returns false if stat classifier returns false for a filter", func(t *testing.T) {
		t.Helper()

		rc := new(MockResultClassifier)
		sc := new(MockStatClassifier)

		rc.On("MatchesFilter", ctx, fix, f1).Return(true, nil)
		rc.On("MatchesFilter", ctx, fix, f2).Return(true, nil)

		sc.On("MatchesFilter", ctx, fix, f3).Return(true, nil)
		sc.On("MatchesFilter", ctx, fix, f4).Return(false, nil)

		matcher := classify.NewFilterMatcher(rc, sc)

		matches, err := matcher.MatchesFilters(ctx, fix, results, stats)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.False(t, matches)
		rc.AssertExpectations(t)
		sc.AssertExpectations(t)
	})

	t.Run("returns error if returned by result classifier", func(t *testing.T) {
		t.Helper()

		rc := new(MockResultClassifier)
		sc := new(MockStatClassifier)

		e := errors.New("error from classifier")

		rc.On("MatchesFilter", ctx, fix, f1).Return(true, nil)
		rc.On("MatchesFilter", ctx, fix, f2).Return(false, e)

		sc.AssertNotCalled(t, "MatchesFilter")

		matcher := classify.NewFilterMatcher(rc, sc)

		_, err := matcher.MatchesFilters(ctx, fix, results, stats)

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, e, err)
		rc.AssertExpectations(t)
		sc.AssertExpectations(t)
	})

	t.Run("returns error if returned by stat classifier", func(t *testing.T) {
		t.Helper()

		rc := new(MockResultClassifier)
		sc := new(MockStatClassifier)

		e := errors.New("error from classifier")

		rc.On("MatchesFilter", ctx, fix, f1).Return(true, nil)
		rc.On("MatchesFilter", ctx, fix, f2).Return(true, nil)

		sc.On("MatchesFilter", ctx, fix, f3).Return(true, nil)
		sc.On("MatchesFilter", ctx, fix, f4).Return(false, e)

		matcher := classify.NewFilterMatcher(rc, sc)

		_, err := matcher.MatchesFilters(ctx, fix, results, stats)

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

func (m *MockResultClassifier) MatchesFilter(ctx context.Context, fix *classify.Fixture, f *trader.ResultFilter) (bool, error) {
	args := m.Called(ctx, fix, f)
	return args.Get(0).(bool), args.Error(1)
}

type MockStatClassifier struct {
	mock.Mock
}

func (m *MockStatClassifier) MatchesFilter(ctx context.Context, fix *classify.Fixture, f *trader.StatFilter) (bool, error) {
	args := m.Called(ctx, fix, f)
	return args.Get(0).(bool), args.Error(1)
}
