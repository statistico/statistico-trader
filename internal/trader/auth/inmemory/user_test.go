package inmemory_test

import (
	"github.com/google/uuid"
	"github.com/statistico/statistico-strategy/internal/trader/auth/inmemory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserService_ByID(t *testing.T) {
	id := uuid.MustParse("f34fcc71-9089-4dd8-a128-25f0725a55d7")

	service := inmemory.NewUserService(id, "test@email.com", "test", "password", "key123")

	t.Run("returns user struct if ID matches expecting expected ID", func(t *testing.T) {
		t.Helper()

		user, err := service.ByID(id)

		if err != nil {
			t.Fatalf("Expected nil, got %+v", err)
		}

		a := assert.New(t)

		a.Equal("f34fcc71-9089-4dd8-a128-25f0725a55d7", user.ID.String())
		a.Equal("test@email.com", user.Email)
		a.Equal("test", user.BetFairUserName)
		a.Equal("password", user.BetFairPassword)
		a.Equal("key123", user.BetFairKey)
	})

	t.Run("returns NotFoundError if ID provided does not match expected ID", func(t *testing.T) {
		t.Helper()

		id := uuid.MustParse("b90a3f79-7c79-4ce8-a0c1-599ad2ab0fd7")

		_, err := service.ByID(id)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "user with ID 'b90a3f79-7c79-4ce8-a0c1-599ad2ab0fd7' does not exist", err.Error())
	})
}
