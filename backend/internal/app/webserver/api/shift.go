package api

//Seems redundant ATM
//Will use models.Shift

//
//type ShiftType int
//
//const (
//	MotorizedPatrolShiftType ShiftType = iota
//	StaticPostShiftType
//	ProactiveOperationShiftType
//	// DailyDutyShiftType could be kitchen duty, HQ duty, etc.
//	DailyDutyShiftType
//)
//
////TODO - consider separating req and resp types where shift will be returned with soldier data(and not just ID)
//
//type Shift struct {
//	StartTimeHour         int       `json:"startTimeHour" validate:"required,min=0,max=24"`
//	StartTimeMinute       int       `json:"startTimeMinute" validate:"required,min=0,max=60"`
//	EndTimeHour           int       `json:"endTimeHour" validate:"required,min=0,max=24"`
//	EndTimeMinute         int       `json:"endTimeMinute" validate:"required,min=0,max=60"`
//	ID                    string    `json:"id" validate:"omitempty"`
//	Name                  string    `json:"name" validate:"required"`
//	Type                  ShiftType `json:"type" validate:"min=0"`
//	CommanderSoldierID    string    `json:"commander" validate:"required"`
//	AdditionalSoldiersIDs []string  `json:"additionalSoldiersIDs"`
//	Description           string    `json:"description" validate:"omitempty,min=1,max=255"`
//	ShiftTemplateID       string    `json:"shiftTemplateId" validate:"omitempty"`
//}
