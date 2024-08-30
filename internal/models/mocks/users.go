package mocks

import (
	"time"

	"github.com/noonacedia/sourcepaste/internal/models"
)

var mockUser = models.User{
	ID:             1,
	Name:           "TestUser",
	Email:          "test@mail.ru",
	HashedPassword: []byte("shitpassword"),
	Created:        time.Now(),
}

type UserModel struct{}

func (u *UserModel) InsertUser(name, email, password string) error {
	return nil
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case mockUser.Email:
		return mockUser.ID, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

func (u *UserModel) Exists(id int) (bool, error) {
	switch id {
	case mockUser.ID:
		return true, nil
	default:
		return false, nil
	}
}
