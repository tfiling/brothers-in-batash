package controllers

import (
	"brothers_in_batash/internal/pkg/store"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

const (
	APIRouteBasePath = "/api/v1"

	RegisterRoute     = "/auth/register"
	LoginRoute        = "/auth/login"
	RefreshTokenRoute = "/auth/refresh"

	CreateShiftRoute  = "/shifts"
	GetShiftRoute     = "/shifts/:id"
	GetAllShiftsRoute = "/shifts"
	UpdateShiftRoute  = "/shifts/:id"
	DeleteShiftRoute  = "/shifts/:id"

	CreateDayScheduleRoute  = "/day-schedules"
	GetDayScheduleRoute     = "/day-schedules/:date"
	GetAllDaySchedulesRoute = "/day-schedules"
	UpdateDayScheduleRoute  = "/day-schedules/:date"
	DeleteDayScheduleRoute  = "/day-schedules/:date"

	CreateSoldierRoute  = "/soldiers"
	GetSoldierRoute     = "/soldiers/:id"
	GetAllSoldiersRoute = "/soldiers"
	UpdateSoldierRoute  = "/soldiers/:id"
	DeleteSoldierRoute  = "/soldiers/:id"
)

type Controller interface {
	RegisterRoutes(router fiber.Router) error
}

func SetupRoutes(v1Router fiber.Router, controllers []Controller) error {
	for _, controller := range controllers {
		if err := controller.RegisterRoutes(v1Router); err != nil {
			return errors.Wrap(err, "failed to register routes")
		}
	}
	return nil
}

func InitControllers() (controllers []Controller, err error) {
	userStore, err := store.NewUserStore()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize user store")
	}

	registrationController, err := NewRegistrationController(userStore)
	controllers = append(controllers, registrationController)
	return
}
