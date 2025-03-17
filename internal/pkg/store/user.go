package store

import (
	"brothers_in_batash/internal/pkg/models"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

//TODO - accept ctx in signatures

type IUserStore interface {
	CreateNewUser(user models.User) error
	FindUserByUsername(username string) ([]models.User, error)
}

type InMemUserStore struct {
	users map[string]models.User
}

func NewUserStore() (*InMemUserStore, error) {
	return &InMemUserStore{users: make(map[string]models.User)}, nil
}

func (us *InMemUserStore) CreateNewUser(user models.User) error {
	if err := validator.New().Struct(user); err != nil {
		return errors.Wrap(err, "user validation failed")
	}
	if _, exists := us.users[user.Username]; exists {
		return errors.New("user already exists")
	}
	us.users[user.Username] = user
	return nil
}

func (us *InMemUserStore) FindUserByUsername(username string) ([]models.User, error) {
	if res, exists := us.users[username]; !exists {
		return []models.User{}, nil
	} else {
		return []models.User{res}, nil
	}
}
