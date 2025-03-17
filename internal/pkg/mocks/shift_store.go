package mocks

import (
	"brothers_in_batash/internal/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockIShiftStore struct {
	mock.Mock
}

func (m *MockIShiftStore) CreateNewShift(shift models.Shift) error {
	args := m.Called(shift)
	return args.Error(0)
}

func (m *MockIShiftStore) FindShiftByID(id string) ([]models.Shift, error) {
	args := m.Called(id)
	return args.Get(0).([]models.Shift), args.Error(1)
}

func (m *MockIShiftStore) FindAllShifts() ([]models.Shift, error) {
	args := m.Called()
	return args.Get(0).([]models.Shift), args.Error(1)
}

func (m *MockIShiftStore) UpdateShift(shift models.Shift) error {
	args := m.Called(shift)
	return args.Error(0)
}

func (m *MockIShiftStore) DeleteShift(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
