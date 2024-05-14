package store

import (
	"brothers_in_batash/internal/pkg/models"
	"time"

	"github.com/pkg/errors"
)

type IDayStore interface {
	CreateNewDaySchedule(day models.DaySchedule) error
	FindDaySchedule(date time.Time) ([]models.DaySchedule, error)
	UpdateDaySchedule(day models.DaySchedule) error
}

type InMemDaySchedStore struct {
	//days maps between a normalized string representation of date to the instance
	days map[string]models.DaySchedule
}

func NewInMemDaySchedStore() (*InMemDaySchedStore, error) {
	return &InMemDaySchedStore{days: make(map[string]models.DaySchedule)}, nil
}

func (s *InMemDaySchedStore) CreateNewDaySchedule(day models.DaySchedule) error {
	if err := day.IsValid(); err != nil {
		return errors.Wrap(err, "Could not insert invalid day schedule")
	}
	s.days[normalizeDate(day.Date)] = day
	return nil
}

func (s *InMemDaySchedStore) FindDaySchedule(date time.Time) ([]models.DaySchedule, error) {
	if day, ok := s.days[normalizeDate(date)]; !ok {
		return []models.DaySchedule{}, nil
	} else {
		return []models.DaySchedule{day}, nil
	}
}

func (s *InMemDaySchedStore) UpdateDaySchedule(day models.DaySchedule) error {
	//TODO - potential bug - users might change day.Date and overwrite the wrong day instance
	if _, ok := s.days[normalizeDate(day.Date)]; !ok {
		return errors.New("day does not exist")
	} else if err := day.IsValid(); err != nil {
		return errors.Wrap(err, "could not update day schedule with an invalid instance")
	}
	s.days[normalizeDate(day.Date)] = day
	return nil
}

func normalizeDate(date time.Time) string {
	return date.Format("2006-01-02")
}
