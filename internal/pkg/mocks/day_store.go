package mocks

import (
	"brothers_in_batash/internal/pkg/models"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockIDayStore struct {
	mock.Mock
}

func (m *MockIDayStore) CreateNewDaySchedule(day models.DaySchedule) error {
	args := m.Called(day)
	return args.Error(0)
}

func (m *MockIDayStore) FindDaySchedule(date time.Time) ([]models.DaySchedule, error) {
	args := m.Called(date)
	return args.Get(0).([]models.DaySchedule), args.Error(1)
}

func (m *MockIDayStore) FindAllDaySchedules() ([]models.DaySchedule, error) {
	args := m.Called()
	return args.Get(0).([]models.DaySchedule), args.Error(1)
}

func (m *MockIDayStore) UpdateDaySchedule(day models.DaySchedule) error {
	args := m.Called(day)
	return args.Error(0)
}

func (m *MockIDayStore) DeleteDaySchedule(date time.Time) error {
	args := m.Called(date)
	return args.Error(0)
}
