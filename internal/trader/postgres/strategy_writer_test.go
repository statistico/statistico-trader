package postgres_test

import (
	"github.com/google/uuid"
	"github.com/statistico/statistico-trader/internal/trader"
	"github.com/statistico/statistico-trader/internal/trader/postgres"
	"github.com/statistico/statistico-trader/internal/trader/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStrategyWriter_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, []string{"strategy", "strategy_result_filter", "strategy_stat_filter"})
	repo := postgres.NewStrategyWriter(conn)

	t.Run("increases tables counts", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		min := float32(1.50)
		max := float32(2.50)
		compIDs := []uint64{8, 12}

		strategyCounts := []struct {
			Strategy      *trader.Strategy
			StrategyCount int8
			FilterCount   int8
		}{
			{
				newStrategy("Strategy One", "Strategy Description", uuid.New(), &min, &max, "MATCH_ODDS", "Home", "BACK", "ACTIVE", "PUBLIC", compIDs),
				1,
				2,
			},
			{
				newStrategy("Strategy Two", "Strategy Description", uuid.New(), &min, nil,"MATCH_ODDS", "Home", "BACK", "ACTIVE", "PUBLIC", compIDs),
				2,
				4,
			},
			{
				newStrategy("Strategy Three", "Strategy Description", uuid.New(), nil, &max, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", compIDs),
				3,
				6,
			},
		}

		for _, sc := range strategyCounts {
			insertStrategy(t, repo, sc.Strategy)

			var strategyCount int8
			var resultCount int8
			var statCount int8

			row := conn.QueryRow("select count(*) from strategy")

			if err := row.Scan(&strategyCount); err != nil {
				t.Errorf("Error when scanning rows returned by the database: %s", err.Error())
			}

			assert.Equal(t, sc.StrategyCount, strategyCount)

			row = conn.QueryRow("select count(*) from strategy_result_filter")

			if err := row.Scan(&resultCount); err != nil {
				t.Errorf("Error when scanning rows returned by the database: %s", err.Error())
			}

			assert.Equal(t, sc.FilterCount, resultCount)

			row = conn.QueryRow("select count(*) from strategy_stat_filter")

			if err := row.Scan(&statCount); err != nil {
				t.Errorf("Error when scanning rows returned by the database: %s", err.Error())
			}

			assert.Equal(t, sc.FilterCount, statCount)
		}
	})

	t.Run("returns a DuplicationError if insert a strategy with a name that exists for user", func(t *testing.T) {
		t.Helper()

		t.Helper()
		defer cleanUp()

		userID := uuid.New()

		stOne := newStrategy("Strategy One", "My Strategy", userID, nil, nil,"MATCH_ODDS", "Home", "BACK", "ACTIVE", "PUBLIC", []uint64{8, 12})
		stTwo := newStrategy("Strategy One", "My Strategy", userID, nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12})

		err := repo.Insert(stOne)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		err = repo.Insert(stTwo)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "Duplication error: Strategy exists with name provided", err.Error())
	})
}

func insertStrategy(t *testing.T, repo trader.StrategyWriter, s *trader.Strategy) {
	if err := repo.Insert(s); err != nil {
		t.Errorf("Error when inserting strategy into the database: %s", err.Error())
	}
}

func newStrategy(
	name string,
	description string,
	userID uuid.UUID,
	min,
	max *float32,
	market,
	runner,
	side,
	status,
	vis string,
	compIDs []uint64,
) *trader.Strategy {
	return &trader.Strategy{
		ID:             uuid.New(),
		Name:           name,
		Description:    description,
		UserID:         userID,
		MarketName:     market,
		RunnerName:     runner,
		MinOdds:        min,
		MaxOdds:        max,
		CompetitionIDs: compIDs,
		Side:           side,
		Visibility:     vis,
		Status:         status,
		ResultFilters: []*trader.ResultFilter{
			{
				Team:   "HOME",
				Result: "WIN",
				Games:  3,
				Venue:  "HOME_AWAY",
			},
			{
				Team:   "AWAY",
				Result: "LOSE",
				Games:  3,
				Venue:  "HOME_AWAY",
			},
		},
		StatFilters: []*trader.StatFilter{
			{
				Stat:    "SHOTS_ON_GOAL",
				Team:    "HOME",
				Action:  "FOR",
				Games:   2,
				Measure: "TOTAL",
				Metric:  "GTE",
				Value:   2,
				Venue:   "AWAY",
			},
			{
				Stat:    "GOALS",
				Team:    "AWAY",
				Action:  "FOR",
				Games:   2,
				Measure: "TOTAL",
				Metric:  "GTE",
				Value:   2,
				Venue:   "AWAY",
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
