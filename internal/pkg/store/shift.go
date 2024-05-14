package store

import (
	"brothers_in_batash/internal/pkg/models"

	"github.com/pkg/errors"
)

type IShiftStore interface {
	CreateNewShift(shift models.Shift) error
	FindShiftByID(id string) (models.Shift, error)
	FindAllShifts() ([]models.Shift, error)
	UpdateShift(shift models.Shift) error
	DeleteShift(id string) error
}

type InMemShiftStore struct {
	shifts map[string]models.Shift
}

func NewShiftStore() (*InMemShiftStore, error) {
	return &InMemShiftStore{shifts: make(map[string]models.Shift)}, nil
}

func (s *InMemShiftStore) CreateNewShift(shift models.Shift) error {
	if err := shift.IsValid(); err != nil {
		return errors.Wrap(err, "shift validation failed")
	}
	if _, exists := s.shifts[shift.ID]; exists {
		return errors.New("shift already exists")
	}
	s.shifts[shift.ID] = shift
	return nil
}

func (s *InMemShiftStore) FindShiftByID(id string) (models.Shift, error) {
	if shift, exists := s.shifts[id]; !exists {
		return models.Shift{}, errors.New("shift not found")
	} else {
		return shift, nil
	}
}

func (s *InMemShiftStore) FindAllShifts() ([]models.Shift, error) {
	shifts := make([]models.Shift, 0, len(s.shifts))
	for _, shift := range s.shifts {
		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (s *InMemShiftStore) UpdateShift(shift models.Shift) error {
	if err := shift.IsValid(); err != nil {
		return errors.Wrap(err, "shift validation failed")
	}
	if _, exists := s.shifts[shift.ID]; !exists {
		return errors.New("shift not found")
	}
	s.shifts[shift.ID] = shift
	return nil
}

func (s *InMemShiftStore) DeleteShift(id string) error {
	if _, exists := s.shifts[id]; !exists {
		return errors.New("shift not found")
	}
	delete(s.shifts, id)
	return nil
}
