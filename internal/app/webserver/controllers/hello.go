package controllers

import (
	"brothers_in_batash/internal/pkg/logging"
	jwtmw "brothers_in_batash/internal/pkg/middleware/jwt"
	"fmt"

	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v5"
)

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

func (c *HelloController) ProtectedHello(ctx *fiber.Ctx) error {
	logging.Info("protected hello", nil)
	token := ctx.Locals(jwtmw.ContextKey).(*jtoken.Token)
	logging.Info(fmt.Sprintf("token: %v", token), nil)
	claims := token.Claims.(jtoken.MapClaims)
	logging.Info(fmt.Sprintf("claims: %v", claims), nil)
	return ctx.SendString("Welcome 👋")
}
