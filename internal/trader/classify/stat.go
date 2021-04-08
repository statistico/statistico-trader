package classify

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader"
	"time"
)

type StatFilterClassifier interface {
	MatchesFilter(ctx context.Context, fix *Fixture, f *trader.StatFilter) (bool, error)
}

type statFilterClassifier struct {
	resultClient statisticodata.ResultClient
}

func (s *statFilterClassifier) MatchesFilter(ctx context.Context, fix *Fixture, f *trader.StatFilter) (bool, error) {
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

	results, err := s.resultClient.ByTeam(ctx, &req)

	if err != nil {
		return false, err
	}

	return statMeetsCriteria(results, teamID, f)
}

func NewStatFilterClassifier(c statisticodata.ResultClient) StatFilterClassifier {
	return &statFilterClassifier{resultClient: c}
}
