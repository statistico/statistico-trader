package trade_test

import (
	"github.com/google/uuid"
	"github.com/statistico/statistico-trader/internal/trader/test"
	"github.com/statistico/statistico-trader/internal/trader/trade"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTradeWriter_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, []string{"trade"})
	writer := trade.NewPostgresWriter(conn)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		tradeCounts := []struct{
			Trade   *trade.Trade
			TradeCount    uint8
		} {
			{
				newTrade(uuid.New(), "SUCCESS"),
				1,
			},
			{
				newTrade(uuid.New(), "IN_PLAY"),
				2,
			},
			{
				newTrade(uuid.New(), "FAIL"),
				3,
			},
		}

		for _, tc := range tradeCounts {
			insertTrade(t, writer, tc.Trade)

			var count uint8

			row := conn.QueryRow("select count(*) from trade")

			if err := row.Scan(&count); err != nil {
				t.Errorf("Error when scanning rows returned by the database: %s", err.Error())
			}

			assert.Equal(t, tc.TradeCount, count)
		}
	})
}

func insertTrade(t *testing.T, w trade.Writer, tr *trade.Trade) {
	if err := w.Insert(tr); err != nil {
		t.Fatalf("Error inserting trade: %s", err.Error())
	}
}

func newTrade(strategyID uuid.UUID, result string) *trade.Trade {
	return &trade.Trade{
		ID:          uuid.New(),
		StrategyID:  strategyID,
		Exchange:    "betfair",
		ExchangeRef: "REF123",
		Market:      "MATCH_ODDS",
		Runner:      "Home",
		Price:       1.90,
		Stake:       100.00,
		EventID:     281781,
		EventDate:   time.Now(),
		Side:        "BACK",
		Result:      result,
		Timestamp:   time.Now(),
	}
}
