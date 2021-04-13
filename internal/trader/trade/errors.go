package trade

import (
	"fmt"
	"github.com/google/uuid"
)

type DuplicationError struct {
	market  string
	runner  string
	eventID uint64
	strategyID uuid.UUID
}

func (d *DuplicationError) Error() string {
	return fmt.Sprintf(
		"trade exists for market %s, runner %s, event %d and strategy %s",
		d.market,
		d.runner,
		d.eventID,
		d.strategyID.String(),
	)
}

type ExchangeError struct {
	err error
}

func (e *ExchangeError) Error() string {
	return fmt.Sprintf("error returned by exchange client: %+v", e.err)
}

type InvalidBalanceError struct {
	market  string
	runner  string
	eventID uint64
	strategyID uuid.UUID
	balance float32
}

func (i *InvalidBalanceError) Error() string {
	return fmt.Sprintf(
		"invalid balance of %.2f when placing trade for market %s, runner %s, event %d and strategy %s",
		i.balance,
		i.market,
		i.runner,
		i.eventID,
		i.strategyID.String(),
	)
}