package store

import (
	"brothers_in_batash/internal/pkg/models"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type ISoldierStore interface {
	CreateNewSoldier(soldier models.Soldier) error
	FindSoldierByID(id string) ([]models.Soldier, error)
	FindAllSoldiers() ([]models.Soldier, error)
	UpdateSoldier(soldier models.Soldier) error
	DeleteSoldier(id string) error
}

type InMemSoldierStore struct {
	soldiers map[string]models.Soldier
}

func NewSoldierStore() (*InMemSoldierStore, error) {
	return &InMemSoldierStore{soldiers: make(map[string]models.Soldier)}, nil
}

func (s *InMemSoldierStore) CreateNewSoldier(soldier models.Soldier) error {
	if err := validator.New().Struct(soldier); err != nil {
		return errors.Wrap(err, "soldier validation failed")
	}
	if _, exists := s.soldiers[soldier.ID]; exists {
		return errors.New("soldier already exists")
	}
	s.soldiers[soldier.ID] = soldier
	return nil
}

func (s *InMemSoldierStore) FindSoldierByID(id string) ([]models.Soldier, error) {
	if soldier, exists := s.soldiers[id]; !exists {
		return []models.Soldier{}, nil
	} else {
		return []models.Soldier{soldier}, nil
	}
}

func (s *InMemSoldierStore) FindAllSoldiers() ([]models.Soldier, error) {
	soldiers := make([]models.Soldier, 0, len(s.soldiers))
	for _, soldier := range s.soldiers {
		soldiers = append(soldiers, soldier)
	}
	return soldiers, nil
}

func (s *InMemSoldierStore) UpdateSoldier(soldier models.Soldier) error {
	if err := validator.New().Struct(soldier); err != nil {
		return errors.Wrap(err, "soldier validation failed")
	}
	if _, exists := s.soldiers[soldier.ID]; !exists {
		return errors.New("soldier not found")
	}
	s.soldiers[soldier.ID] = soldier
	return nil
}

func (s *InMemSoldierStore) DeleteSoldier(id string) error {
	if _, exists := s.soldiers[id]; !exists {
		return errors.New("soldier not found")
	}
	delete(s.soldiers, id)
	return nil
}
