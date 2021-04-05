package classify

import (
	"context"
	"fmt"
	"github.com/statistico/statistico-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
)

type ResultParser interface {
	Parse(ctx context.Context, eventID uint64, market, runner, side string) (Result, error)
}

type resultParser struct {
	resultClient  statisticodata.ResultClient
}

func (r *resultParser) Parse(ctx context.Context, eventID uint64, market, runner, side string) (Result, error) {
	result, err := r.resultClient.ByID(ctx, eventID)

	if err != nil {
		return "", err
	}

	home, away, err := parseGoalScored(result)

	if err != nil {
		return "", err
	}

	res, err := parseResult(market, runner, home, away)

	if err != nil {
		return "", err
	}

	return transformResultForSide(side, res)
}

func parseGoalScored(r *statistico.Result) (uint32, uint32, error) {
	if r.GetStats() == nil {
		return 0, 0, fmt.Errorf("unable to parse match stats for fixture %d", r.Id)
	}

	if r.GetStats().GetHomeScore() == nil {
		return 0, 0, fmt.Errorf("unable to parse home team goals for fixture %d", r.Id)
	}

	if r.GetStats().GetAwayScore() == nil {
		return 0, 0, fmt.Errorf("unable to parse away team goals for fixture %d", r.Id)
	}

	stats := r.GetStats()

	return stats.GetHomeScore().GetValue(), stats.GetAwayScore().GetValue(), nil
}

func parseResult(market, runner string, home, away uint32) (Result, error) {
	switch market {
	case MatchOdds:
		return getMatchOddsResult(market, runner, home, away)
	case OverUnder05:
		return getOverUnderGoalsResult(market, runner, home, away, 0)
	case OverUnder15:
		return getOverUnderGoalsResult(market, runner, home, away, 1)
	case OverUnder25:
		return getOverUnderGoalsResult(market, runner, home, away, 2)
	case OverUnder35:
		return getOverUnderGoalsResult(market, runner, home, away, 3)
	case OverUnder45:
		return getOverUnderGoalsResult(market, runner, home, away, 4)
	default:
		return Fail, fmt.Errorf("market %s is not supported", market)
	}
}

func getMatchOddsResult(market, runner string, home, away uint32) (Result, error) {
	if runner == Home {
		if home > away {
			return Success, nil
		}

		return Fail, nil
	}

	if runner == Away {
		if away > home {
			return Success, nil
		}

		return Fail, nil
	}

	if runner == Draw {
		if home == away {
			return Success, nil
		}

		return Fail, nil
	}

	return Fail, returnRunnerError(market, runner)
}

func getOverUnderGoalsResult(market, runner string, home, away, goals uint32) (Result, error) {
	total := home + away

	if runner[0:4] == Over {
		if total > goals {
			return Success, nil
		}

		return Fail, nil
	}

	if runner[0:5] == Under {
		if total <= goals {
			return Success, nil
		}

		return Fail, nil
	}

	return Fail, returnRunnerError(market, runner)
}

func transformResultForSide(side string, result Result) (Result, error) {
	if side == Back {
		return result, nil
	}

	if side == Lay {
		if result == Success {
			return Fail, nil
		}

		return Success, nil
	}

	return Fail, fmt.Errorf("side %s is not supported", side)
}

func returnRunnerError(market, runner string) error {
	return fmt.Errorf("runner %s not support for market %s", runner, market)
}

func NewResultParser(r statisticodata.ResultClient) ResultParser {
	return &resultParser{resultClient: r}
}
