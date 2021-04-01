package postgres_test

import (
	"github.com/google/uuid"
	"github.com/statistico/statistico-strategy/internal/trader"
	"github.com/statistico/statistico-strategy/internal/trader/postgres"
	"github.com/statistico/statistico-strategy/internal/trader/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTradeReader_Get(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, []string{"trade"})
	writer := postgres.NewTradeWriter(conn)
	reader := postgres.NewTradeReader(conn)

	t.Run("returns a slice of trade.Trade struct", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		strategyID := uuid.New()

		insertTrade(t, writer, newTrade(strategyID, "IN_PLAY"))
		insertTrade(t, writer, newTrade(uuid.New(), "SUCCESS"))
		insertTrade(t, writer, newTrade(strategyID, "FAIL"))
		insertTrade(t, writer, newTrade(strategyID, "IN_PLAY"))

		query := trader.TradeReaderQuery{StrategyID: strategyID}

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
			Query  *trader.TradeReaderQuery
			Count  int
		} {
			{
				&trader.TradeReaderQuery{
					StrategyID: stIdOne,
					Status:     []string{"IN_PLAY", "FAIL"},
				},
				2,
			},
			{
				&trader.TradeReaderQuery{
					StrategyID: stIdTwo,
					Status:     []string{"IN_PLAY", "FAIL"},
				},
				1,
			},
			{
				&trader.TradeReaderQuery{
					StrategyID: stIdOne,
					Status:     []string{"SUCCESS", "FAIL"},
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
