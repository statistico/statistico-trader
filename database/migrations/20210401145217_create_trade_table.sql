-- +goose Up
-- +goose StatementBegin
CREATE TABLE trade (
    id VARCHAR NOT NULL PRIMARY KEY,
    strategy_id VARCHAR NOT NULL,
    exchange VARCHAR NOT NULL,
    exchange_ref VARCHAR NOT NULL,
    market VARCHAR NOT NULL,
    runner VARCHAR NOT NULL,
    price FLOAT NOT NULL,
    stake FLOAT NOT NULL,
    event_id INTEGER NOT NULL,
    event_date INTEGER NOT NULL,
    side VARCHAR NOT NULL,
    result VARCHAR NOT NULL,
    timestamp INTEGER NOT NULL
);

CREATE INDEX on trade (strategy_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trade;
-- +goose StatementEnd
