package strategy_test

import (
	"github.com/google/uuid"
	"github.com/statistico/statistico-trader/internal/trader/strategy"
	"github.com/statistico/statistico-trader/internal/trader/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStrategyReader_Get(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, []string{"strategy", "strategy_result_filter", "strategy_stat_filter"})
	writer := strategy.NewPostgresWriter(conn)
	reader := strategy.NewPostgresReader(conn)

	t.Run("returns a slice of trader.Strategy struct", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		min := float32(1.50)
		max := float32(2.50)

		st := []struct {
			Strategy *strategy.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), &min, &max, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy B", "Second Strategy", uuid.New(), &min, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy C", "Third Strategy", uuid.New(), nil, &max, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		order := "name_asc"

		s, err := reader.Get(&strategy.ReaderQuery{OrderBy: &order})

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
			Strategy *strategy.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy B", "Second Strategy", userID, nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy C", "Third Strategy", userID, nil, nil,"MATCH_ODDS", "Home", "BACK", "ACTIVE", "PUBLIC", []uint64{8, 12}),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		order := "name_asc"

		query := strategy.ReaderQuery{
			OrderBy: &order,
			UserID:  &userID,
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
			Strategy *strategy.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy B", "Second Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy C", "Third Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PRIVATE", []uint64{8, 12}),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		order := "name_asc"
		visibility := "PUBLIC"

		query := strategy.ReaderQuery{
			OrderBy:    &order,
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

	t.Run("strategies can be filtered by Competition ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		st := []struct {
			Strategy *strategy.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy B", "Second Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{5}),
			},
			{
				newStrategy("Strategy C", "Third Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PRIVATE", []uint64{8, 5}),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		strategyCounts := []struct{
			CompetitionID uint64
			Count  int
		} {
			{
				8,
				2,
			},
			{
				5,
				2,
			},
			{
				12,
				1,
			},
			{
				66,
				0,
			},
		}

		for _, sc := range strategyCounts {
			query := strategy.ReaderQuery{CompetitionID: &sc.CompetitionID}

			s, err := reader.Get(&query)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.Equal(t, sc.Count, len(s))
		}
	})

	t.Run("strategies can be filtered by side", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		st := []struct {
			Strategy *strategy.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy B", "Second Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "LAY", "ACTIVE","PUBLIC", []uint64{5}),
			},
			{
				newStrategy("Strategy C", "Third Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PRIVATE", []uint64{8, 5}),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		strategyCounts := []struct{
			Side string
			Count  int
		} {
			{
				"BACK",
				2,
			},
			{
				"LAY",
				1,
			},
		}

		for _, sc := range strategyCounts {
			query := strategy.ReaderQuery{Side: &sc.Side}

			s, err := reader.Get(&query)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.Equal(t, sc.Count, len(s))
		}
	})

	t.Run("strategies can be filtered by market", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		st := []struct {
			Strategy *strategy.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy B", "Second Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Home", "LAY", "ACTIVE","PUBLIC", []uint64{5}),
			},
			{
				newStrategy("Strategy C", "Third Strategy", uuid.New(), nil, nil, "OVER_UNDER_25", "Home", "BACK", "ACTIVE","PRIVATE", []uint64{8, 5}),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		strategyCounts := []struct{
			Market string
			Count  int
		} {
			{
				"MATCH_ODDS",
				2,
			},
			{
				"OVER_UNDER_25",
				1,
			},
			{
				"OVER_UNDER_35",
				0,
			},
		}

		for _, sc := range strategyCounts {
			query := strategy.ReaderQuery{Market: &sc.Market}

			s, err := reader.Get(&query)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.Equal(t, sc.Count, len(s))
		}
	})

	t.Run("strategies can be filtered by runner", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		st := []struct {
			Strategy *strategy.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Away", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy B", "Second Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Away", "LAY", "ACTIVE","PUBLIC", []uint64{5}),
			},
			{
				newStrategy("Strategy C", "Third Strategy", uuid.New(), nil, nil, "OVER_UNDER_25", "Home", "BACK", "ACTIVE","PRIVATE", []uint64{8, 5}),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		strategyCounts := []struct{
			Runner string
			Count  int
		} {
			{
				"Away",
				2,
			},
			{
				"Home",
				1,
			},
			{
				"Draw",
				0,
			},
		}

		for _, sc := range strategyCounts {
			query := strategy.ReaderQuery{Runner: &sc.Runner}

			s, err := reader.Get(&query)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.Equal(t, sc.Count, len(s))
		}
	})

	t.Run("strategies can be filtered by status", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		st := []struct {
			Strategy *strategy.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Away", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy B", "Second Strategy", uuid.New(), nil, nil, "MATCH_ODDS", "Away", "LAY", "ARCHIVED","PUBLIC", []uint64{5}),
			},
			{
				newStrategy("Strategy C", "Third Strategy", uuid.New(), nil, nil, "OVER_UNDER_25", "Home", "BACK", "ACTIVE","PRIVATE", []uint64{8, 5}),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		strategyCounts := []struct{
			Status string
			Count  int
		} {
			{
				"ARCHIVED",
				1,
			},
			{
				"ACTIVE",
				2,
			},
		}

		for _, sc := range strategyCounts {
			query := strategy.ReaderQuery{Status: &sc.Status}

			s, err := reader.Get(&query)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.Equal(t, sc.Count, len(s))
		}
	})

	t.Run("strategies can be filtered by price", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		min := float32(1.50)
		max := float32(3.30)

		st := []struct {
			Strategy *strategy.Strategy
		}{
			{
				newStrategy("Strategy A", "First Strategy", uuid.New(), &min, &max, "MATCH_ODDS", "Away", "BACK", "ACTIVE","PUBLIC", []uint64{8, 12}),
			},
			{
				newStrategy("Strategy B", "Second Strategy", uuid.New(), &min, nil, "MATCH_ODDS", "Away", "LAY", "ARCHIVED","PUBLIC", []uint64{5}),
			},
			{
				newStrategy("Strategy C", "Third Strategy", uuid.New(), &max, nil, "OVER_UNDER_25", "Home", "BACK", "ACTIVE","PRIVATE", []uint64{8, 5}),
			},
		}

		for _, s := range st {
			if err := writer.Insert(s.Strategy); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		strategyCounts := []struct{
			Price float32
			Count  int
		} {
			{
				2.50,
				2,
			},
			{
				5.00,
				2,
			},
			{
				1.25,
				0,
			},
			{
				1.65,
				2,
			},
			{
				3.25,
				2,
			},
		}

		for _, sc := range strategyCounts {
			query := strategy.ReaderQuery{Price: &sc.Price}

			s, err := reader.Get(&query)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.Equal(t, sc.Count, len(s))
		}
	})
}

func assertStrategy(t *testing.T, expected, actual *strategy.Strategy) {
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
