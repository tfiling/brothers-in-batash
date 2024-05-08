package controllers

import (
	"brothers_in_batash/internal/pkg/store"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

const (
	APIRouteBasePath = "/api/v1"
	RegisterRoute    = "/users/register"
)

type Controller interface {
	RegisterRoutes(router fiber.Router) error
}

func SetupRoutes(v1Router fiber.Router) error {
	controllers, err := initControllers()
	if err != nil {
		return errors.Wrap(err, "failed to initialize controllers")
	}
	for _, controller := range controllers {
		if err := controller.RegisterRoutes(v1Router); err != nil {
			return errors.Wrap(err, "failed to register routes")
		}
	}
	return nil
}

func initControllers() (controllers []Controller, err error) {
	userStore, err := store.NewUserStore()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize user store")
	}

	registrationController, err := NewRegistrationController(userStore)
	controllers = append(controllers, registrationController)
	return
}
