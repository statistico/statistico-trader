package strategy

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/statistico/statistico-trader/internal/trader/errors"
)

type PostgresWriter struct {
	connection *sql.DB
}

func (w *PostgresWriter) Insert(s *Strategy) error {
	var exists bool

	err := w.connection.
		QueryRow(`SELECT exists (SELECT id FROM strategy where name = $1 and user_id = $2)`, s.Name, s.UserID.String()).
		Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return &errors.DuplicationError{Message: "Strategy exists with name provided"}
	}

	compIds := make([]int64, len(s.CompetitionIDs))

	for i, c := range s.CompetitionIDs {
		compIds[i] = int64(c)
	}

	builder := queryBuilder(w.connection)

	_, err = builder.
		Insert("strategy").
		Columns(
			"id",
			"name",
			"description",
			"user_id",
			"market",
			"runner",
			"min_odds",
			"max_odds",
			"competition_ids",
			"side",
			"visibility",
			"status",
			"staking_plan",
			"created_at",
			"updated_at",
		).
		Values(
			s.ID.String(),
			s.Name,
			s.Description,
			s.UserID.String(),
			s.MarketName,
			s.RunnerName,
			s.MinOdds,
			s.MaxOdds,
			pq.Array(compIds),
			s.Side,
			s.Visibility,
			s.Status,
			s.StakingPlan,
			s.CreatedAt.Unix(),
			s.UpdatedAt.Unix(),
		).
		Exec()

	err = w.insertResultFilters(s.ID, s.ResultFilters)
	err = w.insertStatFilters(s.ID, s.StatFilters)

	return err
}

func (w *PostgresWriter) insertResultFilters(strategyID uuid.UUID, f []*ResultFilter) error {
	builder := queryBuilder(w.connection)

	for _, filter := range f {
		_, err := builder.
			Insert("strategy_result_filter").
			Columns(
				"strategy_id",
				"team",
				"result",
				"games",
				"venue",
			).
			Values(
				strategyID.String(),
				filter.Team,
				filter.Result,
				filter.Games,
				filter.Venue,
			).
			Exec()

		if err != nil {
			return err
		}
	}

	return nil
}

func (w *PostgresWriter) insertStatFilters(strategyID uuid.UUID, f []*StatFilter) error {
	builder := queryBuilder(w.connection)

	for _, filter := range f {
		_, err := builder.
			Insert("strategy_stat_filter").
			Columns(
				"strategy_id",
				"stat",
				"team",
				"action",
				"measure",
				"metric",
				"games",
				"value",
				"venue",
			).
			Values(
				strategyID.String(),
				filter.Stat,
				filter.Team,
				filter.Action,
				filter.Measure,
				filter.Metric,
				filter.Games,
				filter.Value,
				filter.Venue,
			).
			Exec()

		if err != nil {
			return err
		}
	}

	return nil
}

func NewPostgresWriter(connection *sql.DB) Writer {
	return &PostgresWriter{connection: connection}
}
