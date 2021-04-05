package bootstrap

import (
	"github.com/statistico/statistico-strategy/internal/trader/auth"
	"github.com/statistico/statistico-strategy/internal/trader/auth/config"
)

func (c Container) TokenAuthoriser() auth.TokenAuthoriser {
	aws := c.Config.AWS

	return auth.NewAwsTokenAuthoriser(aws.Region, aws.CognitoUserPoolID, c.Clock, c.Logger)
}

func (c Container) UserService() auth.UserService {
	return config.NewUserService(c.Config)
}
