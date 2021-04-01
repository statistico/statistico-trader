package postgres

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/statistico/statistico-strategy/internal/trader"
	"time"
)

type strategyReader struct {
	connection *sql.DB
}

func (r *strategyReader) Get(q *trader.StrategyReaderQuery) ([]*trader.Strategy, error) {
	query := buildReaderQuery(r.connection, q)

	var id string
	var userID string
	var compIDs []int64
	var created int64
	var updated int64

	rows, err := query.Query()

	if err != nil {
		return []*trader.Strategy{}, err
	}

	defer rows.Close()

	st := []*trader.Strategy{}

	for rows.Next() {
		var s trader.Strategy

		err := rows.Scan(
			&id,
			&s.Name,
			&s.Description,
			&userID,
			&s.MarketName,
			&s.RunnerName,
			&s.MinOdds,
			&s.MaxOdds,
			(*pq.Int64Array)(&compIDs),
			&s.Side,
			&s.Visibility,
			&s.Status,
			&s.StakingPlan,
			&created,
			&updated,
		)

		if err != nil {
			return st, err
		}

		rf, err := r.fetchResultFilters(id)

		if err != nil {
			return st, err
		}

		sf, err := r.fetchStatFilters(id)

		if err != nil {
			return st, err
		}

		ids := make([]uint64, len(compIDs))

		for i, c := range compIDs {
			ids[i] = uint64(c)
		}

		s.ID = uuid.MustParse(id)
		s.UserID = uuid.MustParse(userID)
		s.CompetitionIDs = ids
		s.ResultFilters = rf
		s.StatFilters = sf
		s.CreatedAt = time.Unix(created, 0)
		s.UpdatedAt = time.Unix(updated, 0)

		st = append(st, &s)
	}

	return st, nil
}

func (r *strategyReader) fetchResultFilters(id string) ([]*trader.ResultFilter, error) {
	filters := []*trader.ResultFilter{}

	builder := queryBuilder(r.connection)

	rows, err := builder.
		Select(
			"team",
			"result",
			"games",
			"venue",
		).
		From("strategy_result_filter").
		Where(sq.Eq{"strategy_id": id}).
		Query()

	if err != nil {
		return filters, err
	}

	defer rows.Close()

	for rows.Next() {
		var f trader.ResultFilter

		err := rows.Scan(
			&f.Team,
			&f.Result,
			&f.Games,
			&f.Venue,
		)

		if err != nil {
			return filters, err
		}

		filters = append(filters, &f)
	}

	return filters, nil
}

func (r *strategyReader) fetchStatFilters(id string) ([]*trader.StatFilter, error) {
	filters := []*trader.StatFilter{}

	builder := queryBuilder(r.connection)

	rows, err := builder.
		Select(
			"stat",
			"team",
			"action",
			"measure",
			"metric",
			"games",
			"value",
			"venue",
		).
		From("strategy_stat_filter").
		Where(sq.Eq{"strategy_id": id}).
		Query()

	if err != nil {
		return filters, err
	}

	defer rows.Close()

	for rows.Next() {
		var f trader.StatFilter

		err := rows.Scan(
			&f.Stat,
			&f.Team,
			&f.Action,
			&f.Measure,
			&f.Metric,
			&f.Games,
			&f.Value,
			&f.Venue,
		)

		if err != nil {
			return filters, err
		}

		filters = append(filters, &f)
	}

	return filters, nil
}

func buildReaderQuery(db *sql.DB, q *trader.StrategyReaderQuery) sq.SelectBuilder {
	builder := queryBuilder(db)

	query := builder.Select("strategy.*").From("strategy")

	if q.UserID != nil {
		query = query.Where(sq.Eq{"user_id": q.UserID.String()})
	}

	if q.Market != nil {
		query = query.Where(sq.Eq{"market": *q.Market})
	}

	if q.Runner != nil {
		query = query.Where(sq.Eq{"runner": *q.Runner})
	}

	if q.Price != nil {
		query = query.
			Where(sq.GtOrEq{"max_odds": *q.Price}).
			Where(sq.LtOrEq{"min_odds": *q.Price})
	}

	if q.CompetitionID != nil {
		query = query.Where(" = ANY(competition_ids)?", *q.CompetitionID)
	}

	if q.Visibility != nil {
		query = query.Where(sq.Eq{"visibility": *q.Visibility})
	}

	if q.OrderBy != nil {
		if *q.OrderBy == "name_asc" {
			query.OrderBy("name ASC")
		}

		if *q.OrderBy == "name_desc" {
			query.OrderBy("name DESC")
		}

		if *q.OrderBy == "created_at_asc" {
			query.OrderBy("created_at ASC")
		}

		if *q.OrderBy == "created_at_desc" {
			query.OrderBy("created_at DESC")
		}
	}

	return query
}

func NewStrategyReader(connection *sql.DB) trader.StrategyReader {
	return &strategyReader{connection: connection}
}
