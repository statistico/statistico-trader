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
				newStrategy("Strategy A", "First Strategy", uuid.New(), &min, &max),
			},
			{
				newStrategy("Strategy B", "Second Strategy", uuid.New(), &min, nil),
			},
			{
				newStrategy("Strategy C", "Third Strategy", uuid.New(), nil, &max),
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

		assertStrategy(t, st[0].Strategy, s[0])
		assertStrategy(t, st[1].Strategy, s[1])
		assertStrategy(t, st[2].Strategy, s[2])
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
}
