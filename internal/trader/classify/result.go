package classify

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader"
	"time"
)

type ResultFilterClassifier interface {
	MatchesFilter(ctx context.Context, fix *Fixture, f *trader.ResultFilter) (bool, error)
}

type resultFilterClassifier struct {
	resultClient statisticodata.ResultClient
}

func (r *resultFilterClassifier) MatchesFilter(ctx context.Context, fix *Fixture, f *trader.ResultFilter) (bool, error) {
	teamID, err := parseTeamID(fix, f.Team)

	if err != nil {
		return false, err
	}

	req := statistico.TeamResultRequest{
		TeamId:     teamID,
		Limit:      &wrappers.UInt64Value{Value: uint64(f.Games)},
		DateBefore: &wrappers.StringValue{Value: fix.Date.Format(time.RFC3339)},
		SeasonIds:  []uint64{fix.SeasonID},
		Venue:      &wrappers.StringValue{Value: f.Venue},
	}

	results, err := r.resultClient.ByTeam(ctx, &req)

	if err != nil {
		return false, err
	}

	for _, res := range results {
		if !resultMeetsCriteria(res, teamID, f.Result) {
			return false, nil
		}
	}

	return true, nil
}

func parseTeamID(fix *Fixture, team string) (uint64, error) {
	if team == HomeTeam {
		return fix.HomeTeamID, nil
	}

	if team == AwayTeam {
		return fix.AwayTeamID, nil
	}

	return 0, fmt.Errorf("team enum %s is not supported", team)
}

func resultMeetsCriteria(rs *statistico.Result, teamID uint64, result string) bool {
	switch result {
	case Win:
		return meetsWinCriteria(teamID, rs)
	case Lose:
		return meetsLoseCriteria(teamID, rs)
	case WinDraw:
		return meetsWinDrawCriteria(teamID, rs)
	case LoseDraw:
		return meetsLoseDrawCriteria(teamID, rs)
	case WinLose:
		return meetsWinLoseCriteria(rs)
	default:
		return rs.Stats.GetHomeScore().GetValue() == rs.Stats.GetAwayScore().GetValue()
	}
}

func NewResultFilterClassifier(c statisticodata.ResultClient) ResultFilterClassifier {
	return &resultFilterClassifier{resultClient: c}
}
