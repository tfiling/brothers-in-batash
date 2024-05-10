package controllers_test

import (
	"brothers_in_batash/internal/app/webserver/api"
	"brothers_in_batash/internal/app/webserver/controllers"
	"brothers_in_batash/internal/pkg/store"
	"brothers_in_batash/internal/pkg/test_utils"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegistrationController_NewRegistrationController__error_on_nil_store(t *testing.T) {
	res, err := controllers.NewRegistrationController(nil)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestRegistrationController_NewRegistrationController__success(t *testing.T) {
	userStore, err := store.NewUserStore()
	assert.NoError(t, err)
	assert.NotNil(t, userStore)
	res, err := controllers.NewRegistrationController(userStore)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestRegistrationController_RegisterUser__sad_flows(t *testing.T) {
	//TODO - use a mock instead of the real instance
	testCases := []struct {
		name               string
		body               io.Reader
		expectedStatusCode int
	}{
		{
			"empty body",
			bytes.NewReader([]byte{}),
			http.StatusBadRequest,
		},
		{
			"invalid username",
			test_utils.WrapStructWithReader(t, api.UserRegistrationReqBody{
				Password: "password",
			}),
			http.StatusBadRequest,
		},
		{
			"invalid password",
			test_utils.WrapStructWithReader(t, api.UserRegistrationReqBody{
				Username: "user",
			}),
			http.StatusBadRequest,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//Arrange
			app := fiber.New()

			userStore, err := store.NewUserStore()
			assert.NoError(t, err)
			assert.NotNil(t, userStore)
			controller, err := controllers.NewRegistrationController(userStore)
			assert.NoError(t, err)
			assert.NotNil(t, controller)
			err = controllers.SetupRoutes(app)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, controllers.RegisterRoute, testCase.body)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			//Act
			resp, err := app.Test(req)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)
		})
	}
}

func TestRegistrationController_RegisterUser__user_already_exists(t *testing.T) {
	//Arrange
	app := fiber.New()

	userStore, err := store.NewUserStore()
	assert.NoError(t, err)
	assert.NotNil(t, userStore)
	controller, err := controllers.NewRegistrationController(userStore)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app)
	assert.NoError(t, err)

	newUserBocy := api.UserRegistrationReqBody{
		Username: "user",
		Password: "password",
	}
	req := httptest.NewRequest(fiber.MethodPost, controllers.RegisterRoute, test_utils.WrapStructWithReader(t, newUserBocy))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	secondReq := httptest.NewRequest(fiber.MethodPost, controllers.RegisterRoute, test_utils.WrapStructWithReader(t, newUserBocy))
	secondReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(req)
	resp, err = app.Test(secondReq)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestRegistrationController_RegisterUser__user_successfully_created(t *testing.T) {
	//Arrange
	app := fiber.New()

	userStore, err := store.NewUserStore()
	assert.NoError(t, err)
	assert.NotNil(t, userStore)
	controller, err := controllers.NewRegistrationController(userStore)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app)
	assert.NoError(t, err)

	newUserBocy := api.UserRegistrationReqBody{
		Username: "user",
		Password: "password",
	}
	req := httptest.NewRequest(fiber.MethodPost, controllers.RegisterRoute, test_utils.WrapStructWithReader(t, newUserBocy))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(req)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
