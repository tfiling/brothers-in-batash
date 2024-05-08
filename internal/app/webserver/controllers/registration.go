package controllers

import (
	"brothers_in_batash/internal/app/webserver/api"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationController struct {
	userStore store.IUserStore
}

func NewRegistrationController(userStore store.IUserStore) (*RegistrationController, error) {
	if userStore == nil {
		return nil, errors.New("userStore is nil")
	}
	return &RegistrationController{userStore: userStore}, nil
}

func (c *RegistrationController) RegisterRoutes(router fiber.Router) error {
	router.Post(RegisterRoute, c.registerUser)
	return nil
}

func (c *RegistrationController) registerUser(ctx *fiber.Ctx) error {
	reqBody := api.UserRegistrationReqBody{}
	if err := ctx.BodyParser(&reqBody); err != nil {
		fmt.Printf("Could not parse user registration request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := validator.New().Struct(reqBody); err != nil {
		fmt.Printf("User registration request failed validation: %v", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid registration request body")
	}

	hashedPassword, err := hashPassword(reqBody.Password)
	if err != nil {
		fmt.Printf("Could not hash user pasword: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
	}

	if res, err := c.userStore.FindUserByUsername(reqBody.Username); err != nil {
		fmt.Printf("Could not lookup if user exists: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Could not lookup if user exists")
	} else if len(res) > 0 {
		return ctx.Status(fiber.StatusConflict).SendString("Username already taken")
	}
	newUser := models.User{
		Username:       reqBody.Username,
		HashedPassword: hashedPassword,
	}

	if err := c.userStore.CreateNewUser(newUser); err != nil {
		fmt.Printf("Could not add new user: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Could not add new user")
	}
	return ctx.SendStatus(fiber.StatusCreated)
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}
