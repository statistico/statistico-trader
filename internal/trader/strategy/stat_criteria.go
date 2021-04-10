package strategy

import (
	"fmt"
	"github.com/statistico/statistico-proto/go"
)

func statMeetsCriteria(rs []*statistico.Result, teamID uint64, f *StatFilter) (bool, error) {
	values, err := parseStatValues(rs, teamID, f)

	if err != nil {
		return false, err
	}

	switch f.Measure {
	case Average:
		return meetsAverageCriteria(values, f.Metric, f.Value)
	case Continuous:
		return meetsContinuousCriteria(values, f.Metric, f.Value)
	case Total:
		return meetsTotalCriteria(values, f.Metric, f.Value)
	default:
		return false, fmt.Errorf("metric %s is not supported", f.Metric)
	}
}

func parseStatValues(rs []*statistico.Result, teamID uint64, f *StatFilter) ([]uint32, error) {
	var values []uint32

	for _, res := range rs {
		stats, err := parseTeamStats(res, teamID, f.Action)

		if err != nil {
			return values, err
		}

		val, err := parseStatValue(stats, f.Stat)

		if err != nil {
			return values, err
		}

		values = append(values, val)
	}

	return values, nil
}

func parseTeamStats(res *statistico.Result, teamID uint64, action string) (*statistico.TeamStats, error) {
	if action == ActionFor {
		if res.HomeTeam.Id == teamID {
			if res.HomeTeamStats == nil {
				return nil, fmt.Errorf("no stats available for team %d and result %d", teamID, res.Id)
			}

			return res.HomeTeamStats, nil
		}

		if res.AwayTeamStats == nil {
			return nil, fmt.Errorf("no stats available for team %d and result %d", teamID, res.Id)
		}

		return res.AwayTeamStats, nil
	}

	if action == ActionAgainst {
		if res.HomeTeam.Id == teamID {
			if res.AwayTeamStats == nil {
				return nil, fmt.Errorf("no stats available for team %d and result %d", teamID, res.Id)
			}

			return res.AwayTeamStats, nil
		}

		if res.HomeTeamStats == nil {
			return nil, fmt.Errorf("no stats available for team %d and result %d", teamID, res.Id)
		}

		return res.HomeTeamStats, nil
	}

	return nil, fmt.Errorf("action %s is not supported", action)
}

func parseStatValue(s *statistico.TeamStats, stat string) (uint32, error) {
	switch stat {
	case Goals:
		return s.Goals.GetValue(), nil
	case ShotsOnGoal:
		return s.ShotsOnGoal.GetValue(), nil
	default:
		return 0, fmt.Errorf("stat %s is not supported", stat)
	}
}

func meetsAverageCriteria(values []uint32, metric string, value float32) (bool, error) {
	var val uint32

	for _, v := range values {
		val += v
	}

	calc := float32(val) / float32(len(values))

	if metric == Gte {
		return (float32(int(calc*100)) / 100) >= value, nil
	}

	if metric == Lte {
		return (float32(int(calc*100)) / 100) <= value, nil
	}

	return false, fmt.Errorf("metric %s is not supported", metric)
}

func meetsContinuousCriteria(values []uint32, metric string, value float32) (bool, error) {
	for _, v := range values {
		if metric == Gte {
			if float32(v) < value {
				return false, nil
			}

			continue
		}

		if metric == Lte {
			if float32(v) > value {
				return false, nil
			}

			continue
		}

		return false, fmt.Errorf("metric %s is not supported", metric)
	}

	return true, nil
}

func meetsTotalCriteria(values []uint32, metric string, value float32) (bool, error) {
	var val uint32

	for _, v := range values {
		val += v
	}

	calc := float32(val)

	if metric == Gte {
		return (float32(int(calc*100)) / 100) >= value, nil
	}

	if metric == Lte {
		return (float32(int(calc*100)) / 100) <= value, nil
	}

	return false, fmt.Errorf("metric %s is not supported", metric)
}
