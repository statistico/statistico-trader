-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategy (
    id VARCHAR NOT NULL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    user_id VARCHAR NOT NULL,
    market VARCHAR NOT NULL,
    runner VARCHAR NOT NULL,
    min_odds FLOAT,
    max_odds FLOAT,
    competition_ids INTEGER[] NOT NULL,
    side VARCHAR NOT NULL,
    visibility VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    staking_plan JSON NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE strategy_result_filter (
    strategy_id VARCHAR NOT NULL,
    team VARCHAR NOT NULL,
    result VARCHAR NOT NULL,
    games SMALLINT NOT NULL,
    venue VARCHAR NOT NULL,
    CONSTRAINT fk_strategy
        FOREIGN KEY(strategy_id)
            REFERENCES strategy(id)
            ON DELETE CASCADE
);

CREATE TABLE strategy_stat_filter (
    strategy_id VARCHAR NOT NULL,
    stat VARCHAR NOT NULL,
    team VARCHAR NOT NULL,
    action VARCHAR NOT NULL,
    measure VARCHAR NOT NULL,
    metric VARCHAR NOT NULL,
    games SMALLINT NOT NULL,
    value SMALLINT NOT NULL,
    venue VARCHAR NOT NULL,
    CONSTRAINT fk_strategy
        FOREIGN KEY(strategy_id)
            REFERENCES strategy(id)
            ON DELETE CASCADE
);

CREATE INDEX ON strategy (user_id);
CREATE INDEX ON strategy_result_filter (strategy_id);
CREATE INDEX ON strategy_stat_filter (strategy_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE strategy;
DROP TABLE strategy_result_filter;
DROP TABLE strategy_stat_filter;
-- +goose StatementEnd
