package controllers

import (
	"brothers_in_batash/internal/app/webserver/api"
	"brothers_in_batash/internal/pkg/logging"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"errors"

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
		logging.Info("Could not parse user registration request body", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	if err := validator.New().Struct(reqBody); err != nil {
		logging.Info("User registration request failed validation", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	hashedPassword, err := hashPassword(reqBody.Password)
	if err != nil {
		logging.Info("Could not hash user password", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if res, err := c.userStore.FindUserByUsername(reqBody.Username); err != nil {
		logging.Info("Could not lookup if user exists", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	} else if len(res) > 0 {
		//TODO - reconsider this, security-wise
		logging.Debug("Registration attempt with a take username", []logging.LogProp{{"username", reqBody.Username}})
		return ctx.Status(fiber.StatusConflict).SendString("Username already taken")
	}
	newUser := models.User{
		Username:       reqBody.Username,
		HashedPassword: hashedPassword,
	}

	if err := c.userStore.CreateNewUser(newUser); err != nil {
		logging.Warning(err, "error on writing new user to DB", nil)
		return ctx.SendStatus(fiber.StatusInternalServerError)
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
