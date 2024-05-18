package controllers

import (
	"brothers_in_batash/internal/pkg/config"
	"brothers_in_batash/internal/pkg/middleware/jwt"
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

type storeInstancesContainer struct {
	dayStore     store.IDayStore
	shiftStore   store.IShiftStore
	soldierStore store.ISoldierStore
	userStore    store.IUserStore
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
	storeInstances, err := initStoreInstances()
	if err != nil {
		return
	}

	authMiddleware := jwt.NewAuthMiddleware(config.JWTSecret)

	registrationController, err := NewRegistrationController(storeInstances.userStore)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize registration controller")
	}
	controllers = append(controllers, registrationController)

	dayScheduleController, err := NewDayScheduleController(storeInstances.dayStore, authMiddleware)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize day schedule controller")
	}
	controllers = append(controllers, dayScheduleController)

	shiftController, err := NewShiftController(storeInstances.shiftStore, authMiddleware)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize shift controller")
	}
	controllers = append(controllers, shiftController)

	soldierController, err := NewSoldierController(storeInstances.soldierStore, authMiddleware)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize shift controller")
	}
	controllers = append(controllers, soldierController)

	return
}

func initStoreInstances() (storeInstancesContainer, error) {
	userStore, err := store.NewUserStore()
	if err != nil {
		return storeInstancesContainer{}, errors.Wrap(err, "failed to initialize user store")
	}

	daySchedStore, err := store.NewInMemDaySchedStore()
	if err != nil {
		return storeInstancesContainer{}, errors.Wrap(err, "failed to initialize day schedule store")
	}

	shiftStore, err := store.NewShiftStore()
	if err != nil {
		return storeInstancesContainer{}, errors.Wrap(err, "failed to initialize shift store")
	}

	soldierStore, err := store.NewSoldierStore()
	if err != nil {
		return storeInstancesContainer{}, errors.Wrap(err, "failed to initialize soldier store")
	}
	return storeInstancesContainer{
		dayStore:     daySchedStore,
		shiftStore:   shiftStore,
		soldierStore: soldierStore,
		userStore:    userStore,
	}, nil
}
