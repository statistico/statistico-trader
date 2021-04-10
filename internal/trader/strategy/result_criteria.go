package strategy

import "github.com/statistico/statistico-proto/go"

func meetsWinCriteria(teamID uint64, rs *statistico.Result) bool {
	stats := rs.GetStats()

	if stats.GetHomeScore().GetValue() == stats.GetAwayScore().GetValue() {
		return false
	}

	if rs.HomeTeam.Id == teamID {
		return stats.GetHomeScore().GetValue() > stats.GetAwayScore().GetValue()
	}

	return stats.GetAwayScore().GetValue() > stats.GetHomeScore().GetValue()
}

func meetsLoseCriteria(teamID uint64, rs *statistico.Result) bool {
	stats := rs.GetStats()

	if stats.GetHomeScore().GetValue() == stats.GetAwayScore().GetValue() {
		return false
	}

	if rs.HomeTeam.Id == teamID {
		return stats.GetHomeScore().GetValue() < stats.GetAwayScore().GetValue()
	}

	return stats.GetAwayScore().GetValue() < stats.GetHomeScore().GetValue()
}

func meetsWinDrawCriteria(teamID uint64, rs *statistico.Result) bool {
	stats := rs.GetStats()

	if stats.GetHomeScore().GetValue() == stats.GetAwayScore().GetValue() {
		return true
	}

	if rs.HomeTeam.Id == teamID {
		return stats.GetHomeScore().GetValue() >= stats.GetAwayScore().GetValue()
	}

	return stats.GetAwayScore().GetValue() >= stats.GetHomeScore().GetValue()
}

func meetsLoseDrawCriteria(teamID uint64, rs *statistico.Result) bool {
	stats := rs.GetStats()

	if stats.GetHomeScore().GetValue() == stats.GetAwayScore().GetValue() {
		return true
	}

	if rs.HomeTeam.Id == teamID {
		return stats.GetHomeScore().GetValue() <= stats.GetAwayScore().GetValue()
	}

	return stats.GetAwayScore().GetValue() <= stats.GetHomeScore().GetValue()
}

func meetsWinLoseCriteria(rs *statistico.Result) bool {
	stats := rs.GetStats()

	return stats.GetHomeScore().GetValue() != stats.GetAwayScore().GetValue()
}
