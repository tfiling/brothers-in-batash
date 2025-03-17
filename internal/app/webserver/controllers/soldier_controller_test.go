package controllers_test

import (
	"brothers_in_batash/internal/app/webserver/controllers"
	"brothers_in_batash/internal/pkg/mocks"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/test_utils"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSoldierController_NewSoldierController__error_on_nil_store(t *testing.T) {
	// Act
	controller, err := controllers.NewSoldierController(nil, test_utils.AlwaysAllowedJWTMiddleware)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestSoldierController_NewSoldierController__error_on_nil_auth_middleware(t *testing.T) {
	// Act
	controller, err := controllers.NewSoldierController(&mocks.MockISoldierStore{}, nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestSoldierController_NewSoldierController__success(t *testing.T) {
	// Arrange
	soldierStore := &mocks.MockISoldierStore{}

	// Act
	controller, err := controllers.NewSoldierController(soldierStore, test_utils.AlwaysAllowedJWTMiddleware)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, controller)
}

func TestSoldierController_CreateSoldier__invalid_request_body(t *testing.T) {
	// Arrange
	app := fiber.New()
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewSoldierController(soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPost, controllers.CreateSoldierRoute, test_utils.WrapStructWithReader(t, "invalid"))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	soldierStore.AssertExpectations(t)
}

func TestSoldierController_CreateSoldier__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewSoldierController(soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	soldier := models.Soldier{ID: "1", FirstName: "John", LastName: "Doe"}
	soldierStore.On("CreateNewSoldier", soldier).Return(nil)
	req := httptest.NewRequest(fiber.MethodPost, controllers.CreateSoldierRoute, test_utils.WrapStructWithReader(t, soldier))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	soldierStore.AssertExpectations(t)
}

func TestSoldierController_GetSoldier__not_found(t *testing.T) {
	// Arrange
	app := fiber.New()
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewSoldierController(soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	soldierID := "1"
	soldierStore.On("FindSoldierByID", soldierID).Return([]models.Soldier{}, nil)
	req := httptest.NewRequest(fiber.MethodGet, "/soldiers/"+soldierID, nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	soldierStore.AssertExpectations(t)
}

func TestSoldierController_GetSoldier__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewSoldierController(soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	soldierID := "1"
	soldier := models.Soldier{ID: soldierID, FirstName: "John", LastName: "Doe"}
	soldierStore.On("FindSoldierByID", soldierID).Return([]models.Soldier{soldier}, nil)
	req := httptest.NewRequest(fiber.MethodGet, "/soldiers/"+soldierID, nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var respSoldier models.Soldier
	err = json.NewDecoder(resp.Body).Decode(&respSoldier)
	assert.NoError(t, err)
	assert.Equal(t, soldier, respSoldier)
	soldierStore.AssertExpectations(t)
}

func TestSoldierController_GetAllSoldiers__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewSoldierController(soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	soldiers := []models.Soldier{
		{ID: "1", FirstName: "John", LastName: "Doe"},
		{ID: "2", FirstName: "Jane", LastName: "Smith"},
	}
	soldierStore.On("FindAllSoldiers").Return(soldiers, nil)
	req := httptest.NewRequest(fiber.MethodGet, controllers.GetAllSoldiersRoute, nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var respSoldiers []models.Soldier
	err = json.NewDecoder(resp.Body).Decode(&respSoldiers)
	assert.NoError(t, err)
	assert.Equal(t, soldiers, respSoldiers)
	soldierStore.AssertExpectations(t)
}

func TestSoldierController_UpdateSoldier__invalid_request_body(t *testing.T) {
	// Arrange
	app := fiber.New()
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewSoldierController(soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	soldierID := "1"
	req := httptest.NewRequest(fiber.MethodPut, "/soldiers/"+soldierID, test_utils.WrapStructWithReader(t, "invalid"))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	soldierStore.AssertExpectations(t)
}

func TestSoldierController_UpdateSoldier__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewSoldierController(soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	soldierID := "1"
	soldier := models.Soldier{ID: soldierID, FirstName: "John", LastName: "Doe"}
	soldierStore.On("UpdateSoldier", mock.AnythingOfType("models.Soldier")).Return(nil)
	req := httptest.NewRequest(fiber.MethodPut, "/soldiers/"+soldierID, test_utils.WrapStructWithReader(t, soldier))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	soldierStore.AssertExpectations(t)
}

func TestSoldierController_DeleteSoldier__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewSoldierController(soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	soldierID := "1"
	soldierStore.On("DeleteSoldier", soldierID).Return(nil)
	req := httptest.NewRequest(fiber.MethodDelete, "/soldiers/"+soldierID, nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	soldierStore.AssertExpectations(t)
}
