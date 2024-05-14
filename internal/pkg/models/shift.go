package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type ShiftType int

const (
	MotorizedPatrolShiftType ShiftType = iota
	StaticPostShiftType
	ProactiveOperationShiftType
	// DailyDutyShiftType could be kitchen duty, HQ duty, etc.
	DailyDutyShiftType
)

type TimeOfDay struct {
	Hour   int `json:"hour" validate:"min=0,max=23"`
	Minute int `json:"minute" validate:"min=0,max=59"`
}

type Shift struct {
	ID                 string    `json:"id" validate:"required"`
	Name               string    `json:"name" validate:"required"`
	Type               ShiftType `json:"type" validate:"min=0"`
	StartTime          TimeOfDay `json:"startTime" validate:"required"`
	EndTime            TimeOfDay `json:"endTime" validate:"required"`
	Commander          Soldier   `json:"commander" validate:"required"`
	AdditionalSoldiers []Soldier `json:"additionalSoldiers"`
	Description        string    `json:"description"`
}

type DaySchedule struct {
	Date   time.Time `json:"date" validate:"required"`
	Shifts []Shift   `json:"shifts" validate:"required,min=1"`
}

func (s Shift) IsValid() error {
	if err := validator.New().Struct(s); err != nil {
		return errors.Wrap(err, "shift failed validation")
	}
	if s.EndTime.Before(s.StartTime) {
		return errors.New("shift end time must be greater than start time")
	}
	return nil
}

func (d DaySchedule) IsValid() error {
	if err := validator.New().Struct(d); err != nil {
		return errors.Wrap(err, "day schedule failed validation")
	}
	for _, shift := range d.Shifts {
		if err := shift.IsValid(); err != nil {
			return errors.Wrap(err, "day schedule member failed validation")
		}
	}
	return nil
}

func (t TimeOfDay) Before(other TimeOfDay) bool {
	if t.Hour < other.Hour {
		return true
	}
	if t.Hour == other.Hour && t.Minute < other.Minute {
		return true
	}
	return false
}
