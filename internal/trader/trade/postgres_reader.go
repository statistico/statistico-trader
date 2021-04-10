package trade

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"time"
)

type PostgresReader struct {
	connection *sql.DB
}

func (r *PostgresReader) Get(q *ReaderQuery) ([]*Trade, error) {
	query := buildTradeQuery(r.connection, q)

	trades := []*Trade{}

	rows, err := query.Query()

	if err != nil {
		return trades, err
	}

	var id string
	var strategyID string
	var eventDate int64
	var timestamp int64

	for rows.Next() {
		var tr Trade

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

func (r *PostgresReader) Exists(market, runner string, eventID uint64, strategyID uuid.UUID) (bool, error) {
	return true, nil
}

func buildTradeQuery(db *sql.DB, q *ReaderQuery) sq.SelectBuilder {
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

func queryBuilder(c *sql.DB) sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(c)
}

func NewPostgresReader(connection *sql.DB) Reader {
	return &PostgresReader{connection: connection}
}
