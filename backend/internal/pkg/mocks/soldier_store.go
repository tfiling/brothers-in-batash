package mocks

import (
	"brothers_in_batash/internal/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockISoldierStore struct {
	mock.Mock
}

func (m *MockISoldierStore) CreateNewSoldier(soldier models.Soldier) error {
	args := m.Called(soldier)
	return args.Error(0)
}

func (m *MockISoldierStore) FindSoldierByID(id string) ([]models.Soldier, error) {
	args := m.Called(id)
	return args.Get(0).([]models.Soldier), args.Error(1)
}

func (m *MockISoldierStore) FindAllSoldiers() ([]models.Soldier, error) {
	args := m.Called()
	return args.Get(0).([]models.Soldier), args.Error(1)
}

func (m *MockISoldierStore) UpdateSoldier(soldier models.Soldier) error {
	args := m.Called(soldier)
	return args.Error(0)
}

func (m *MockISoldierStore) DeleteSoldier(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
