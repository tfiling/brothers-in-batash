package utils

import "github.com/google/uuid"

func NewEntityID() string {
	//ATM, an entity ID is a simple UUID string. Should be good enough for now.
	return uuid.New().String()
}
