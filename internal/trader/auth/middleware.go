package auth

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenAuthoriser interface {
	Authorise(ctx context.Context) (context.Context, error)
}

type awsTokenAuthoriser struct {
	region     string
	userPoolID string
	clock      jwt.Clock
	logger     *logrus.Logger
}

func (a *awsTokenAuthoriser) Authorise(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", a.region, a.userPoolID)

	set, err := jwk.Fetch(context.Background(), fmt.Sprintf(url+"/%s", ".well-known/jwks.json"))

	if err != nil {
		a.logger.Errorf("error fetching key set from aws: %s", err.Error())

		return nil, status.Error(codes.Internal, "internal server error")
	}

	parsed, err := jwt.Parse(
		[]byte(token),
		jwt.WithKeySet(set),
		jwt.WithValidate(true),
		jwt.WithIssuer(url),
		jwt.WithClock(a.clock),
	)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	id, ok := parsed.Get("username")

	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	newCtx := context.WithValue(ctx, "userID", id)

	return newCtx, nil
}

func NewAwsTokenAuthoriser(region, userPoolID string, clock jwt.Clock, logger *logrus.Logger) TokenAuthoriser {
	return &awsTokenAuthoriser{
		region:     region,
		userPoolID: userPoolID,
		clock:      clock,
		logger:     logger,
	}
}
