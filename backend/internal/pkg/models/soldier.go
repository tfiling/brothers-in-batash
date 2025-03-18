package models

type Soldier struct {
	ID             string          `json:"id" validate:"required"`
	FirstName      string          `json:"firstName" validate:"required,alpha"`
	MiddleName     string          `json:"middleName" validate:"omitempty,alpha"`
	LastName       string          `json:"lastName" validate:"required,alpha"`
	PersonalNumber string          `json:"personalNumber" validate:"required,numeric,len=7"`
	Position       SoldierPosition `json:"position" validate:"required"`
	Roles          []SoldierRole   `json:"roles" validate:"min=1,dive"`
}

// SoldierRole is the "Pakal"s of the soldier. Admins will be able to edit and add roles.
// Also used for managing driving qualifications and commanding positions
type SoldierRole struct {
	ID          string `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,alpha"`
	Description string `json:"description" validate:"omitempty,max=255"`
}

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
