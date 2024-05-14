package models

type SoldierPosition int

const (
	PlatoonCommanderPosition SoldierPosition = iota
	VicePlatoonCommanderPosition
	SquadCommanderPosition
	SquadLieutenantPosition
	CommanderPosition
	RegularSoldierPosition
)

func (sp SoldierPosition) String() string {
	return [...]string{
		"Platoon Commander",
		"Vice Platoon Commander",
		"Squad Commander",
		"Squad Lieutenant",
		"Commander",
		"RegularSoldier",
	}[sp]
}

type SoldierRole int

type Soldier struct {
	ID             string          `json:"id" validate:"required"`
	FirstName      string          `json:"firstName" validate:"required,alpha"`
	MiddleName     string          `json:"middleName" validate:"omitempty,alpha"`
	LastName       string          `json:"lastName" validate:"required,alpha"`
	PersonalNumber string          `json:"personalNumber" validate:"required,numeric,len=7"`
	Position       SoldierPosition `json:"position" validate:"required"`
}
