package config

import (
	"github.com/google/uuid"
	"github.com/statistico/statistico-strategy/internal/trader/auth"
	"github.com/statistico/statistico-strategy/internal/trader/bootstrap"
)

type userService struct {
	config  *bootstrap.Config
}

func (u *userService) ByID(userID uuid.UUID) (*auth.User, error) {
	if userID != u.config.User.ID {
		return nil, &auth.NotFoundError{UserID: userID}
	}

	user := u.config.User

	return &auth.User{
		ID:              user.ID,
		Email:           user.Email,
		BetFairUserName: user.BetFairUserName,
		BetFairPassword: user.BetFairPassword,
		BetFairKey:      user.BetFairKey,
	}, nil
}

func NewUserService(c *bootstrap.Config) auth.UserService {
	return &userService{config: c}
}