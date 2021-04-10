package trade

import (
	"database/sql"
)

type PostgresWriter struct {
	connection *sql.DB
}

func (w *PostgresWriter) Insert(t *Trade) error {
	builder := queryBuilder(w.connection)

	_, err := builder.
		Insert("trade").
		Columns(
			"id",
			"exchange",
			"exchange_ref",
			"strategy_id",
			"market",
			"runner",
			"price",
			"stake",
			"event_id",
			"event_date",
			"side",
			"result",
			"timestamp",
		).
		Values(
			t.ID.String(),
			t.Exchange,
			t.ExchangeRef,
			t.StrategyID.String(),
			t.Market,
			t.Runner,
			t.Price,
			t.Stake,
			t.EventID,
			t.EventDate.Unix(),
			t.Side,
			t.Result,
			t.Timestamp.Unix(),
		).Exec()

	if err != nil {
		return err
	}

	return nil
}

func NewPostgresWriter(connection *sql.DB) Writer {
	return &PostgresWriter{connection: connection}
}
