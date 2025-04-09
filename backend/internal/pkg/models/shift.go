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

type ShiftTime struct {
	StartTime TimeOfDay     `json:"startTime" validate:"required"`
	Duration  time.Duration `json:"duration" validate:"required,min=60000000000"`
}

// PersonnelRequirement is a container for any personnel requirements/ constraints for a specific shift
type PersonnelRequirement struct {
	// SoldierRoleToCount maps between specific SoldierRole and the minimum number of soldiers with that Role that are required for the shift
	SoldierRoleToCount map[string]int
}

// ShiftTemplate is used to describe the repetitive tasks. Based on ShiftTemplate, Shift instances are being created.
type ShiftTemplate struct {
	ID                   string               `json:"id" validate:"required"`
	Name                 string               `json:"name" validate:"required"`
	Description          string               `json:"description" validate:"required"`
	PersonnelRequirement PersonnelRequirement `json:"personnelRequirement" validate:"required"`
	// DaysOfOccurrences maps from a weekday to start and end times of a shift.
	//An empty slice would indicate this shift will not happen on that weekday.
	DaysOfOccurrences map[time.Weekday][]ShiftTime `json:"dayOfWeek" validate:"required"`
}

// Shift describes a specific shift, in a specific time and date.
// Shift could be created based on a ShiftTemplate(will provide constraints), or out of scratch.
type Shift struct {
	ID                 string    `json:"id" validate:"required"`
	StartTime          time.Time `json:"startTime" validate:"required"`
	EndTime            time.Time `json:"endTime" validate:"required"`
	Name               string    `json:"name" validate:"required"`
	Type               ShiftType `json:"type" validate:"min=0"`
	Commander          Soldier   `json:"commander" validate:"required"`
	AdditionalSoldiers []Soldier `json:"additionalSoldiers" validate:"dive"`
	Description        string    `json:"description" validate:"omitempty,min=1,max=255"`
	ShiftTemplateID    string    `json:"shiftTemplateId" validate:"omitempty"`
}

type DaySchedule struct {
	Date   time.Time `json:"date" validate:"required"`
	Shifts []Shift   `json:"shifts" validate:"required,min=1"`
}

func (s Shift) IsValid() error {
	if err := validator.New().Struct(s); err != nil {
		return errors.Wrap(err, "shift failed validation")
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
