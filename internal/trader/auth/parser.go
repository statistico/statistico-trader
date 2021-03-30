package auth

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type Token struct {
	UserID  string
}

type TokenParser interface {
	Parse(token string) (*Token, error)
}

type awsTokenParser struct {
	region  string
	userPoolID string
	clock jwt.Clock
}

func (a *awsTokenParser) Parse(token string) (*Token, error) {
	url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", a.region, a.userPoolID)

	set, err := jwk.Fetch(context.Background(), fmt.Sprintf(url + "/%s", ".well-known/jwks.json"))

	if err != nil {
		return nil, &KeySetParseError{err: err}
	}

	parsed, err := jwt.Parse(
		[]byte(token),
		jwt.WithKeySet(set),
		jwt.WithValidate(true),
		jwt.WithIssuer(url),
		jwt.WithClock(a.clock),
	)

	if err != nil {
		return nil, &TokenValidationError{err: err}
	}

	id, ok := parsed.Get("username")

	if !ok {
		return nil, &InvalidFormatError{field: "username"}
	}

	return &Token{
		UserID:  fmt.Sprintf("%v", id),
	}, nil
}

func NewAwsTokenParser(region, userPoolID string, clock jwt.Clock) TokenParser {
	return &awsTokenParser{
		region:     region,
		userPoolID: userPoolID,
		clock:      clock,
	}
}
