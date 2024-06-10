package store

import (
	"brothers_in_batash/internal/pkg/models"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

//TODO - accept ctx in signatures

type IShiftTemplateStore interface {
	CreateNewShiftTemplate(template models.ShiftTemplate) error
	FindShiftTemplateByID(id string) ([]models.ShiftTemplate, error)
	FindAllShiftsTemplate() ([]models.ShiftTemplate, error)
	UpdateShiftTemplate(template models.ShiftTemplate) error
	DeleteShiftTemplate(id string) error
}

type InMemShiftTemplateStore struct {
	shiftTemplates map[string]models.ShiftTemplate
}

func NewShiftTemplateStore() (*InMemShiftTemplateStore, error) {
	return &InMemShiftTemplateStore{shiftTemplates: make(map[string]models.ShiftTemplate)}, nil
}

func (s *InMemShiftTemplateStore) CreateNewShiftTemplate(template models.ShiftTemplate) error {
	if err := validator.New().Struct(template); err != nil {
		return errors.Wrap(err, "shift template validation failed")
	}
	if _, exists := s.shiftTemplates[template.ID]; exists {
		return errors.New("shift template already exists")
	}
	s.shiftTemplates[template.ID] = template
	return nil
}

func (s *InMemShiftTemplateStore) FindShiftTemplateByID(id string) ([]models.ShiftTemplate, error) {
	if template, exists := s.shiftTemplates[id]; !exists {
		return []models.ShiftTemplate{}, nil
	} else {
		return []models.ShiftTemplate{template}, nil
	}
}

func (s *InMemShiftTemplateStore) FindAllShiftsTemplate() ([]models.ShiftTemplate, error) {
	templates := make([]models.ShiftTemplate, 0, len(s.shiftTemplates))
	for _, template := range s.shiftTemplates {
		templates = append(templates, template)
	}
	return templates, nil
}

func (s *InMemShiftTemplateStore) UpdateShiftTemplate(template models.ShiftTemplate) error {
	if err := validator.New().Struct(template); err != nil {
		return errors.Wrap(err, "shift template validation failed")
	}
	if _, exists := s.shiftTemplates[template.ID]; !exists {
		return errors.New("shift template not found")
	}
	s.shiftTemplates[template.ID] = template
	return nil
}

func (s *InMemShiftTemplateStore) DeleteShiftTemplate(id string) error {
	if _, exists := s.shiftTemplates[id]; !exists {
		return errors.New("shift template not found")
	}
	delete(s.shiftTemplates, id)
	return nil
}
