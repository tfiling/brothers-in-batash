package main

import (
	"brothers_in_batash/internal/app/webserver/controllers"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("started WS")

	app := fiber.New()
	apiGroup := app.Group(controllers.APIRouteBasePath)
	helloController, err := controllers.NewHelloController(apiGroup)
	if err != nil {
		fmt.Printf("error creating hello controller: %v", err)
		return
	}
	if err = helloController.SetupRoutes(); err != nil {
		fmt.Printf("error setting up routes: %v", err)
		return
	}

	err = app.Listen(":3000")
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
