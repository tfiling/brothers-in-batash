package mocks

import (
	"brothers_in_batash/internal/pkg/models"

	"github.com/stretchr/testify/mock"
)

// MockIShiftTemplateStore is a mock of IShiftTemplateStore interface.
type MockIShiftTemplateStore struct {
	mock.Mock
}

// CreateNewShiftTemplate mocks base method.
func (m *MockIShiftTemplateStore) CreateNewShiftTemplate(template models.ShiftTemplate) error {
	args := m.Called(template)
	return args.Error(0)
}

// FindShiftTemplateByID mocks base method.
func (m *MockIShiftTemplateStore) FindShiftTemplateByID(id string) ([]models.ShiftTemplate, error) {
	args := m.Called(id)
	return args.Get(0).([]models.ShiftTemplate), args.Error(1)
}

// FindAllShiftsTemplate mocks base method.
func (m *MockIShiftTemplateStore) FindAllShiftsTemplate() ([]models.ShiftTemplate, error) {
	args := m.Called()
	return args.Get(0).([]models.ShiftTemplate), args.Error(1)
}

// UpdateShiftTemplate mocks base method.
func (m *MockIShiftTemplateStore) UpdateShiftTemplate(template models.ShiftTemplate) error {
	args := m.Called(template)
	return args.Error(0)
}

// DeleteShiftTemplate mocks base method.
func (m *MockIShiftTemplateStore) DeleteShiftTemplate(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
