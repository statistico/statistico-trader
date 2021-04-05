package auth

import (
	"fmt"
	"github.com/google/uuid"
)

type NotFoundError struct {
	UserID  uuid.UUID
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("user with ID '%s' does not exist", n.UserID.String())
}
