package main

import (
	"brothers_in_batash/internal/app/webserver/controllers"
	"brothers_in_batash/internal/pkg/logging"

	"github.com/gofiber/fiber/v2"
)

func main() {
	logging.Info("started WS", nil)

	app := fiber.New()
	apiGroup := app.Group(controllers.APIRouteBasePath)
	APIControllers, err := controllers.InitControllers()
	if err != nil {
		logging.Panic(err, "error setting up controllers", nil)
	}
	if err := controllers.SetupRoutes(apiGroup, APIControllers); err != nil {
		logging.Panic(err, "error setting up routes", nil)
	}

	if err := app.Listen(":3000"); err != nil {
		logging.Error(err, "Error starting server", nil)
	}
}
