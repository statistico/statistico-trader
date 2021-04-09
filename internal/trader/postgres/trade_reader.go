package postgres

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/statistico/statistico-trader/internal/trader"
	"time"
)

type TradeReader struct {
	connection *sql.DB
}

func (t *TradeReader) Get(q *trader.TradeReaderQuery) ([]*trader.Trade, error) {
	query := buildTradeQuery(t.connection, q)

	trades := []*trader.Trade{}

	rows, err := query.Query()

	if err != nil {
		return trades, err
	}

	var id string
	var strategyID string
	var eventDate int64
	var timestamp int64

	for rows.Next() {
		var tr trader.Trade

		err := rows.Scan(
			&id,
			&strategyID,
			&tr.Exchange,
			&tr.ExchangeRef,
			&tr.Market,
			&tr.Runner,
			&tr.Price,
			&tr.Stake,
			&tr.EventID,
			&eventDate,
			&tr.Side,
			&tr.Result,
			&timestamp,
		)

		if err != nil {
			return trades, err
		}

		tr.ID = uuid.MustParse(id)
		tr.StrategyID = uuid.MustParse(strategyID)
		tr.EventDate = time.Unix(eventDate, 0)
		tr.Timestamp = time.Unix(timestamp, 0)

		trades = append(trades, &tr)
	}

	return trades, nil
}

func buildTradeQuery(db *sql.DB, q *trader.TradeReaderQuery) sq.SelectBuilder {
	builder := queryBuilder(db)

	query := builder.
		Select("trade.*").
		From("trade").
		Where(sq.Eq{"strategy_id": q.StrategyID.String()})

	if len(q.Result) > 0 {
		query = query.Where(sq.Eq{"result": q.Result})
	}

	return query
}

func NewTradeReader(connection *sql.DB) trader.TradeReader {
	return &TradeReader{connection: connection}
}
