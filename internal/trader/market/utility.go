package market

import (
	"fmt"
	"github.com/statistico/statistico-proto/go"
)

func transformQueryToTrade(q *Query) *Trade {
	return &Trade{
		MarketName:    q.MarketName,
		RunnerName:    q.RunnerName,
		RunnerPrice:   q.RunnerPrice,
		EventId:       q.EventId,
		CompetitionId: q.CompetitionId,
		SeasonId:      q.SeasonId,
		EventDate:     q.EventDate,
		Side:          q.Side,
	}
}

func parseTradeResult(q *Query, r *statistico.Result) (string, error) {
	home, away, err := parseGoalScored(r)

	if err != nil {
		return "", err
	}

	result := Success

	switch q.MarketName {
	case MatchOdds:
		result, err = getMatchOddsResult(q, home, away)
		break
	case OverUnder05:
		result, err = getOverUnderGoalsResult(q, home, away, 0)
		break
	case OverUnder15:
		result, err = getOverUnderGoalsResult(q, home, away, 1)
		break
	case OverUnder25:
		result, err = getOverUnderGoalsResult(q, home, away, 2)
		break
	case OverUnder35:
		result, err = getOverUnderGoalsResult(q, home, away, 3)
		break
	case OverUnder45:
		result, err = getOverUnderGoalsResult(q, home, away, 4)
		break
	default:
		return Fail, fmt.Errorf("market %s is not supported", q.MarketName)
	}

	if err != nil {
		return "", err
	}

	return transformResultForSide(q.Side, result)
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

func getMatchOddsResult(q *Query, home, away uint32) (string, error) {
	if q.RunnerName == Home {
		if home > away {
			return Success, nil
		}

		return Fail, nil
	}

	if q.RunnerName == Away {
		if away > home {
			return Success, nil
		}

		return Fail, nil
	}

	if q.RunnerName == Draw {
		if home == away {
			return Success, nil
		}

		return Fail, nil
	}

	return Fail, returnRunnerError(q.MarketName, q.RunnerName)
}

func getOverUnderGoalsResult(q *Query, home, away, goals uint32) (string, error) {
	total := home + away

	if q.RunnerName[0:4] == Over {
		if total > goals {
			return Success, nil
		}

		return Fail, nil
	}

	if q.RunnerName[0:5] == Under {
		if total <= goals {
			return Success, nil
		}

		return Fail, nil
	}

	return Fail, returnRunnerError(q.MarketName, q.RunnerName)
}

func transformResultForSide(side, result string) (string, error) {
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
