package strategy

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type postgresReader struct {
	connection *sql.DB
}

func (r *postgresReader) Get(q *ReaderQuery) ([]*Strategy, error) {
	query := buildReaderQuery(r.connection, q)

	var id string
	var userID string
	var compIDs []int64
	var created int64
	var updated int64

	rows, err := query.Query()

	if err != nil {
		return []*Strategy{}, err
	}

	defer rows.Close()

	st := []*Strategy{}

	for rows.Next() {
		var s Strategy

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

func (r *postgresReader) fetchResultFilters(id string) ([]*ResultFilter, error) {
	filters := []*ResultFilter{}

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
		var f ResultFilter

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

func (r *postgresReader) fetchStatFilters(id string) ([]*StatFilter, error) {
	filters := []*StatFilter{}

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
		var f StatFilter

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

func buildReaderQuery(db *sql.DB, q *ReaderQuery) sq.SelectBuilder {
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
			Where("(min_odds <= ? OR min_odds IS NULL)", *q.Price).
			Where("(max_odds >= ? OR max_odds IS NULL)", *q.Price)
	}

	if q.CompetitionID != nil {
		query = query.Where("? = ANY(competition_ids)", *q.CompetitionID)
	}

	if q.Side != nil {
		query = query.Where(sq.Eq{"side": *q.Side})
	}

	if q.Status != nil {
		query = query.Where(sq.Eq{"status": *q.Status})
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

func queryBuilder(c *sql.DB) sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(c)
}

func NewPostgresReader(connection *sql.DB) Reader {
	return &postgresReader{connection: connection}
}
