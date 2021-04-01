package postgres

import (
	"database/sql"
	"github.com/statistico/statistico-strategy/internal/trader"
)

type TradeWriter struct {
	connection *sql.DB
}

func (w *TradeWriter) Insert(t *trader.Trade) error {
	builder := queryBuilder(w.connection)

	_, err := builder.
		Insert("trade").
		Columns(
			"id",
			"exchange_ref",
			"strategy_id",
			"market",
			"runner",
			"price",
			"event_id",
			"event_date",
			"side",
			"exchange",
			"result",
			"timestamp",
		).
		Values(
			t.ID.String(),
			t.ExchangeRef,
			t.StrategyID.String(),
			t.Market,
			t.Runner,
			t.Price,
			t.EventID,
			t.EventDate.Unix(),
			t.Side,
			t.Exchange,
			t.Result,
			t.Timestamp.Unix(),
		).Exec()

	if err != nil {
		return err
	}

	return nil
}
