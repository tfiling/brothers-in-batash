package controllers_test

import (
	"brothers_in_batash/internal/app/webserver/api"
	"brothers_in_batash/internal/app/webserver/controllers"
	jwtmw "brothers_in_batash/internal/pkg/middleware/jwt"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"brothers_in_batash/internal/pkg/test_utils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type IUserStoreMock struct {
	mock.Mock
}

func (m *IUserStoreMock) CreateNewUser(user models.User) error {
	return m.Called(user).Error(0)
}

func (m *IUserStoreMock) FindUserByUsername(username string) ([]models.User, error) {
	args := m.Called(username)
	return args.Get(0).([]models.User), args.Error(1)
}

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
	//TODO: remove skip
	t.Skip("ATM I canceled the user registration flow")
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
			err = controllers.SetupRoutes(app, []controllers.Controller{controller})
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, controllers.RegisterRoute, testCase.body)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			//Act
			resp, err := app.Test(req, test_utils.TestTimeout)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)
		})
	}
}

func TestRegistrationController_RegisterUser__user_already_exists(t *testing.T) {
	//TODO: remove skip
	t.Skip("ATM I canceled the user registration flow")

	//Arrange
	app := fiber.New()

	userStore, err := store.NewUserStore()
	assert.NoError(t, err)
	assert.NotNil(t, userStore)
	controller, err := controllers.NewRegistrationController(userStore)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app, []controllers.Controller{controller})
	assert.NoError(t, err)

	newUserBody := api.UserRegistrationReqBody{
		Username: "user",
		Password: "password",
	}
	req := httptest.NewRequest(fiber.MethodPost, controllers.RegisterRoute, test_utils.WrapStructWithReader(t, newUserBody))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	secondReq := httptest.NewRequest(fiber.MethodPost, controllers.RegisterRoute, test_utils.WrapStructWithReader(t, newUserBody))
	secondReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(req, test_utils.TestTimeout)
	resp, err = app.Test(secondReq, test_utils.TestTimeout)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestRegistrationController_RegisterUser__user_successfully_created(t *testing.T) {
	//TODO: remove skip
	t.Skip("ATM I canceled the user registration flow")

	//Arrange
	app := fiber.New()

	userStore, err := store.NewUserStore()
	assert.NoError(t, err)
	assert.NotNil(t, userStore)
	controller, err := controllers.NewRegistrationController(userStore)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app, []controllers.Controller{controller})
	assert.NoError(t, err)

	newUserBody := api.UserRegistrationReqBody{
		Username: "user",
		Password: "password",
	}
	req := httptest.NewRequest(fiber.MethodPost, controllers.RegisterRoute, test_utils.WrapStructWithReader(t, newUserBody))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestRegistrationController_LoginUser__sad_flows(t *testing.T) {
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
			test_utils.WrapStructWithReader(t, api.UserLoginReqBody{
				Password: "password",
			}),
			http.StatusBadRequest,
		},
		{
			"invalid password",
			test_utils.WrapStructWithReader(t, api.UserLoginReqBody{
				Username: "user",
			}),
			http.StatusBadRequest,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//Arrange
			app := fiber.New()

			userStoreMock := &IUserStoreMock{}
			controller, err := controllers.NewRegistrationController(userStoreMock)
			assert.NoError(t, err)
			assert.NotNil(t, controller)
			err = controllers.SetupRoutes(app, []controllers.Controller{controller})
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, controllers.LoginRoute, testCase.body)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			//Act
			resp, err := app.Test(req, test_utils.TestTimeout)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)
			userStoreMock.AssertExpectations(t)
		})
	}
}

func TestRegistrationController_LoginUser__non_existing_user(t *testing.T) {
	//Arrange
	app := fiber.New()

	username := "user"
	userStoreMock := &IUserStoreMock{}
	userStoreMock.On("FindUserByUsername", username).Return([]models.User{}, nil)
	controller, err := controllers.NewRegistrationController(userStoreMock)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app, []controllers.Controller{controller})
	assert.NoError(t, err)

	loginBody := api.UserLoginReqBody{
		Username: username,
		Password: "password",
	}
	req := httptest.NewRequest(fiber.MethodPost, controllers.LoginRoute, test_utils.WrapStructWithReader(t, loginBody))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	userStoreMock.AssertExpectations(t)
}

func TestRegistrationController_LoginUser__wrong_password(t *testing.T) {
	//Arrange
	app := fiber.New()

	username := "user"
	userStoreMock := &IUserStoreMock{}
	userStoreMock.On("FindUserByUsername", username).Return([]models.User{
		{
			Username:       username,
			HashedPassword: []byte("you will never steal my secrets!"),
		},
	}, nil)
	controller, err := controllers.NewRegistrationController(userStoreMock)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app, []controllers.Controller{controller})
	assert.NoError(t, err)

	loginBody := api.UserLoginReqBody{
		Username: username,
		Password: "wrong_password",
	}
	loginReq := httptest.NewRequest(fiber.MethodPost, controllers.LoginRoute, test_utils.WrapStructWithReader(t, loginBody))
	loginReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(loginReq, test_utils.TestTimeout)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	userStoreMock.AssertExpectations(t)
}

func TestRegistrationController_LoginUser__success(t *testing.T) {
	//Arrange
	app := fiber.New()

	username := "user"
	pass := "password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	assert.NoError(t, err)

	userStoreMock := &IUserStoreMock{}
	userStoreMock.On("FindUserByUsername", username).Return([]models.User{
		{
			Username:       username,
			HashedPassword: hashedPassword,
		},
	}, nil)
	controller, err := controllers.NewRegistrationController(userStoreMock)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app, []controllers.Controller{controller})
	assert.NoError(t, err)

	loginBody := api.UserLoginReqBody{
		Username: username,
		Password: pass,
	}
	loginReq := httptest.NewRequest(fiber.MethodPost, controllers.LoginRoute, test_utils.WrapStructWithReader(t, loginBody))
	loginReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(loginReq, test_utils.TestTimeout)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	userStoreMock.AssertExpectations(t)
}

func TestRegistrationController_RefreshToken__sad_flows(t *testing.T) {
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
			"invalid request body",
			test_utils.WrapStructWithReader(t, api.RefreshTokenReqBody{}),
			http.StatusBadRequest,
		},
		{
			"invalid refresh token",
			test_utils.WrapStructWithReader(t, api.RefreshTokenReqBody{
				RefreshToken: "invalid_token",
			}),
			http.StatusUnauthorized,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//Arrange
			app := fiber.New()

			userStoreMock := &IUserStoreMock{}
			controller, err := controllers.NewRegistrationController(userStoreMock)
			assert.NoError(t, err)
			assert.NotNil(t, controller)
			err = controllers.SetupRoutes(app, []controllers.Controller{controller})
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, controllers.RefreshTokenRoute, testCase.body)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			//Act
			resp, err := app.Test(req, test_utils.TestTimeout)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)
			userStoreMock.AssertExpectations(t)
		})
	}
}

func TestRegistrationController_RefreshToken__wrong_claims_type(t *testing.T) {
	//Arrange
	app := fiber.New()

	username := "user"
	wrongClaims := jtoken.MapClaims{
		"myCustomWrongUserID":     username,
		"myCustomWrongExpiration": time.Now().Add(time.Hour).Unix(),
		"myCustomWrongRole":       "admin",
	}
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, wrongClaims)
	signedToken, err := token.SignedString([]byte(jwtmw.SigningSecret))
	assert.NoError(t, err)

	userStoreMock := &IUserStoreMock{}
	controller, err := controllers.NewRegistrationController(userStoreMock)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app, []controllers.Controller{controller})
	assert.NoError(t, err)

	refreshTokenBody := api.RefreshTokenReqBody{
		RefreshToken: signedToken,
	}
	refreshTokenReq := httptest.NewRequest(fiber.MethodPost, controllers.RefreshTokenRoute, test_utils.WrapStructWithReader(t, refreshTokenBody))
	refreshTokenReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(refreshTokenReq, test_utils.TestTimeout)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	userStoreMock.AssertExpectations(t)
}

func TestRegistrationController_RefreshToken__success(t *testing.T) {
	//Arrange
	app := fiber.New()

	username := "user"
	refreshToken, err := jwtmw.GenerateToken(username, jwtmw.RefreshTokenExpiration)
	assert.NoError(t, err)

	userStoreMock := &IUserStoreMock{}
	controller, err := controllers.NewRegistrationController(userStoreMock)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app, []controllers.Controller{controller})
	assert.NoError(t, err)

	refreshTokenBody := api.RefreshTokenReqBody{
		RefreshToken: refreshToken,
	}
	refreshTokenReq := httptest.NewRequest(fiber.MethodPost, controllers.RefreshTokenRoute, test_utils.WrapStructWithReader(t, refreshTokenBody))
	refreshTokenReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(refreshTokenReq, test_utils.TestTimeout)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	userStoreMock.AssertExpectations(t)

	var respBody api.UserLoginRespBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, respBody.Token)
	assert.Equal(t, refreshToken, respBody.RefreshToken)
}

func TestRegistrationController_LogoutUser__invalid_body(t *testing.T) {
	//Arrange
	app := fiber.New()

	userStoreMock := &IUserStoreMock{}
	controller, err := controllers.NewRegistrationController(userStoreMock)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app, []controllers.Controller{controller})
	assert.NoError(t, err)

	req := httptest.NewRequest(fiber.MethodPost, controllers.LogoutRoute, bytes.NewReader([]byte{}))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	userStoreMock.AssertExpectations(t)
}

func TestRegistrationController_LogoutUser__missing_token(t *testing.T) {
	//Arrange
	app := fiber.New()

	userStoreMock := &IUserStoreMock{}
	controller, err := controllers.NewRegistrationController(userStoreMock)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app, []controllers.Controller{controller})
	assert.NoError(t, err)

	logoutBody := api.LogoutReqBody{
		Token: "",
	}
	req := httptest.NewRequest(fiber.MethodPost, controllers.LogoutRoute, test_utils.WrapStructWithReader(t, logoutBody))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	userStoreMock.AssertExpectations(t)
}

func TestRegistrationController_LogoutUser__invalid_token(t *testing.T) {
	//Arrange
	app := fiber.New()

	userStoreMock := &IUserStoreMock{}
	controller, err := controllers.NewRegistrationController(userStoreMock)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app, []controllers.Controller{controller})
	assert.NoError(t, err)

	logoutBody := api.LogoutReqBody{
		Token: "invalid-token",
	}
	req := httptest.NewRequest(fiber.MethodPost, controllers.LogoutRoute, test_utils.WrapStructWithReader(t, logoutBody))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	userStoreMock.AssertExpectations(t)
}

func TestRegistrationController_LogoutUser__success(t *testing.T) {
	//Arrange
	app := fiber.New()

	username := "user"
	token, err := jwtmw.GenerateToken(username, jwtmw.TokenExpiration)
	assert.NoError(t, err)

	userStoreMock := &IUserStoreMock{}
	controller, err := controllers.NewRegistrationController(userStoreMock)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controllers.SetupRoutes(app, []controllers.Controller{controller})
	assert.NoError(t, err)

	logoutBody := api.LogoutReqBody{
		Token: token,
	}
	req := httptest.NewRequest(fiber.MethodPost, controllers.LogoutRoute, test_utils.WrapStructWithReader(t, logoutBody))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	//Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	userStoreMock.AssertExpectations(t)
}
