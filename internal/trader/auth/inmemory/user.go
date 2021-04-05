package inmemory

import (
	"github.com/google/uuid"
	"github.com/statistico/statistico-strategy/internal/trader/auth"
)

type userService struct {
	UserID   string
	UserEmail string
	BetFairUserName string
	BetFairPassword string
	BetFairKey string
}

func (u *userService) ByID(userID uuid.UUID) (*auth.User, error) {
	if userID.String() != u.UserID {
		return nil, &auth.NotFoundError{UserID: userID}
	}

	return &auth.User{
		ID:              uuid.MustParse(u.UserID),
		Email:           u.UserEmail,
		BetFairUserName: u.BetFairUserName,
		BetFairPassword: u.BetFairPassword,
		BetFairKey:      u.BetFairKey,
	}, nil
}

func NewUserService(userID , email, username, password, key string) auth.UserService {
	return &userService{
		UserID:          userID,
		UserEmail:       email,
		BetFairUserName: username,
		BetFairPassword: password,
		BetFairKey:      key,
	}
}
