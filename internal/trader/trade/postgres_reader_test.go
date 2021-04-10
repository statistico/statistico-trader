package trade_test

import (
	"github.com/google/uuid"
	"github.com/statistico/statistico-trader/internal/trader/test"
	"github.com/statistico/statistico-trader/internal/trader/trade"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTradeReader_Get(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, []string{"trade"})
	writer := trade.NewPostgresWriter(conn)
	reader := trade.NewPostgresReader(conn)

	t.Run("returns a slice of trade.Trade struct", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		strategyID := uuid.New()

		insertTrade(t, writer, newTrade(strategyID, "IN_PLAY"))
		insertTrade(t, writer, newTrade(uuid.New(), "SUCCESS"))
		insertTrade(t, writer, newTrade(strategyID, "FAIL"))
		insertTrade(t, writer, newTrade(strategyID, "IN_PLAY"))

		query := trade.ReaderQuery{StrategyID: strategyID}

		trades, err := reader.Get(&query)

		if err != nil {
			t.Fatalf("Expected nil, got %+v", err)
		}

		assert.Equal(t, 3, len(trades))
	})

	t.Run("trades can be filtered by status", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		stIdOne := uuid.New()
		stIdTwo := uuid.New()

		insertTrade(t, writer, newTrade(stIdOne, "IN_PLAY"))
		insertTrade(t, writer, newTrade(stIdTwo, "SUCCESS"))
		insertTrade(t, writer, newTrade(stIdOne, "FAIL"))
		insertTrade(t, writer, newTrade(stIdTwo, "IN_PLAY"))
		insertTrade(t, writer, newTrade(stIdOne, "SUCCESS"))

		tradeCounts := []struct{
			Query  *trade.ReaderQuery
			Count  int
		} {
			{
				&trade.ReaderQuery{
					StrategyID: stIdOne,
					Result:     []string{"IN_PLAY", "FAIL"},
				},
				2,
			},
			{
				&trade.ReaderQuery{
					StrategyID: stIdTwo,
					Result:     []string{"IN_PLAY", "FAIL"},
				},
				1,
			},
			{
				&trade.ReaderQuery{
					StrategyID: stIdOne,
					Result:     []string{"SUCCESS", "FAIL"},
				},
				2,
			},
		}

		for _, tc := range tradeCounts {
			trades, err := reader.Get(tc.Query)

			if err != nil {
				t.Fatalf("Expected nil, got %+v", err)
			}

			assert.Equal(t, tc.Count, len(trades))
		}
	})
}
