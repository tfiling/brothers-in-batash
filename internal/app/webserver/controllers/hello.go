package controllers

import "github.com/gofiber/fiber/v2"

const APIRouteBasePath = "/api/v1"

type Controller interface {
	SetupRoutes() error
}

type HelloController struct {
	router fiber.Router
}

func NewHelloController(router fiber.Router) (*HelloController, error) {
	return &HelloController{router: router}, nil
}

func (c *HelloController) SetupRoutes() error {
	c.router.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	return nil
}
