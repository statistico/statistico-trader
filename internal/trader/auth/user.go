package auth

import "github.com/google/uuid"

// UserService provides an interface to fetch User resources from an abstract data source
type UserService interface {
	ByID(userID uuid.UUID) (*User, error)
}

type User struct {
	ID   uuid.UUID
	Email string
	BetFairUserName string
	BetFairPassword string
	BetFairKey string
}
