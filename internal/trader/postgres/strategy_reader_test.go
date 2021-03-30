package postgres_test

import (
	"github.com/google/uuid"
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/statistico/statistico-strategy/internal/trader/postgres"
	"github.com/statistico/statistico-strategy/internal/trader/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStrategyReader_Get(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, []string{"strategy", "strategy_result_filter", "strategy_stat_filter"})
	writer := postgres.NewStrategyWriter(conn)
	reader := postgres.NewStrategyReader(conn)

	t.Run("returns a slice of trader.Strategy struct", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		min := float32(1.50)
		max := float32(2.50)

		st := []struct {
			Strategy *trader.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), &min, &max, "PUBLIC"),
			},
			{
				newStrategy("Strategy B", "Second Strategy", uuid.New(), &min, nil, "PUBLIC"),
			},
			{
				newStrategy("Strategy C", "Third Strategy", uuid.New(), nil, &max, "PUBLIC"),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		order := "name_asc"

		s, err := reader.Get(&trader.StrategyReaderQuery{OrderBy: &order})

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, 3, len(s))
		assertStrategy(t, st[0].Strategy, s[0])
		assertStrategy(t, st[1].Strategy, s[1])
		assertStrategy(t, st[2].Strategy, s[2])
	})

	t.Run("strategies can be filtered by User ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		userID := uuid.New()

		st := []struct {
			Strategy *trader.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), nil, nil, "PUBLIC"),
			},
			{
				newStrategy("Strategy B", "Second Strategy", userID, nil, nil, "PUBLIC"),
			},
			{
				newStrategy("Strategy C", "Third Strategy", userID, nil, nil, "PUBLIC"),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		order := "name_asc"

		query := trader.StrategyReaderQuery{
			OrderBy: &order,
			UserID: &userID,
		}

		s, err := reader.Get(&query)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, 2, len(s))
		assertStrategy(t, st[1].Strategy, s[0])
		assertStrategy(t, st[2].Strategy, s[1])
	})

	t.Run("strategies can be filtered by visibility", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		st := []struct {
			Strategy *trader.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), nil, nil, "PUBLIC"),
			},
			{
				newStrategy("Strategy B", "Second Strategy", uuid.New(), nil, nil, "PUBLIC"),
			},
			{
				newStrategy("Strategy C", "Third Strategy", uuid.New(), nil, nil, "PRIVATE"),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		order := "name_asc"
		visibility := "PUBLIC"

		query := trader.StrategyReaderQuery{
			OrderBy: &order,
			Visibility: &visibility,
		}

		s, err := reader.Get(&query)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, 2, len(s))
		assertStrategy(t, st[0].Strategy, s[0])
		assertStrategy(t, st[1].Strategy, s[1])
	})
}

func assertStrategy(t *testing.T, expected, actual *trader.Strategy) {
	a := assert.New(t)

	a.Equal(expected.ID, actual.ID)
	a.Equal(expected.Name, actual.Name)
	a.Equal(expected.Description, actual.Description)
	a.Equal(expected.UserID, actual.UserID)
	a.Equal(expected.MarketName, actual.MarketName)
	a.Equal(expected.RunnerName, actual.RunnerName)
	a.Equal(expected.MinOdds, actual.MinOdds)
	a.Equal(expected.MaxOdds, actual.MaxOdds)
	a.Equal(expected.CompetitionIDs, actual.CompetitionIDs)
	a.Equal(expected.Side, actual.Side)
	a.Equal(expected.Visibility, actual.Visibility)
	a.Equal(expected.Status, actual.Status)
	a.Equal(expected.StakingPlan, actual.StakingPlan)
	a.Equal(expected.ResultFilters, actual.ResultFilters)
	a.Equal(expected.StatFilters, actual.StatFilters)
	a.Equal(expected.CreatedAt.Unix(), actual.CreatedAt.Unix())
	a.Equal(expected.UpdatedAt.Unix(), actual.UpdatedAt.Unix())
}
