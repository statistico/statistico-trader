package bootstrap

import "github.com/statistico/statistico-strategy/internal/trader/auth"

func (c Container) TokenAuthoriser() auth.TokenAuthoriser {
	aws := c.Config.AWS

	return auth.NewAwsTokenAuthoriser(aws.Region, aws.CognitoUserPoolID, c.Clock, c.Logger)
}
