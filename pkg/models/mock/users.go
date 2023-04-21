package mock

import (
	"time"

	"github.com/lipandr/Snippetbox/pkg/models"
)

var mockUser = &models.User{
	ID:      1,
	Name:    "Andre",
	Email:   "andre@andre.com",
	Created: time.Now(),
	Active:  true,
}

type UserModel struct{}

func (m *UserModel) Insert(_, email, _ string) error {
	switch email {
	case "andre@andre.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, _ string) (int, error) {
	switch email {
	case "andre@andre.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
