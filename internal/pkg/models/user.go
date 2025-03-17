package models

type User struct {
	Username       string `validate:"ascii,min=4,max=100"`
	HashedPassword []byte `validate:"ascii,min=4,max=100"`
	SoldierID      string `validate:"omitempty"`
}
