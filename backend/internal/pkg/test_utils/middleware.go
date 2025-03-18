package test_utils

import "github.com/gofiber/fiber/v2"

func AlwaysAllowedJWTMiddleware(ctx *fiber.Ctx) error {
	return ctx.Next()
}
