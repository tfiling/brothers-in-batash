package controllers_test

import (
	"brothers_in_batash/internal/app/webserver/controllers"
	"brothers_in_batash/internal/pkg/mocks"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"brothers_in_batash/internal/pkg/test_utils"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	shiftID = "123"
	commanderID
	testShiftName = "Test Shift"
)

var testCommander = models.Soldier{
	ID:             commanderID,
	FirstName:      "Gal",
	LastName:       "Tfilin",
	PersonalNumber: "1234567",
	Position:       models.CommanderPosition,
	Roles: []models.SoldierRole{{
		ID:   "1",
		Name: "Commander",
	}},
}

var testShiftModel = models.Shift{
	ID:                 shiftID,
	Name:               testShiftName,
	Type:               models.MotorizedPatrolShiftType,
	StartTime:          time.Date(2025, time.April, 9, 0, 0, 0, 0, time.UTC),
	EndTime:            time.Date(2025, time.April, 9, 1, 0, 0, 0, time.UTC),
	Commander:          testCommander,
	AdditionalSoldiers: nil,
}

func TestShiftController_NewShiftController__sad_flows(t *testing.T) {
	testCases := []struct {
		shiftStore     store.IShiftStore
		soldierStore   store.ISoldierStore
		authMiddleware fiber.Handler
		name           string
	}{
		{
			shiftStore:     nil,
			soldierStore:   &mocks.MockISoldierStore{},
			authMiddleware: test_utils.AlwaysAllowedJWTMiddleware,
			name:           "nil shift store",
		},
		{
			shiftStore:     &mocks.MockIShiftStore{},
			soldierStore:   nil,
			authMiddleware: test_utils.AlwaysAllowedJWTMiddleware,
			name:           "nil soldier store",
		},
		{
			shiftStore:     &mocks.MockIShiftStore{},
			soldierStore:   &mocks.MockISoldierStore{},
			authMiddleware: nil,
			name:           "nil auth middleware",
		},
	}
	for _, testCase := range testCases {
		// Act
		controller, err := controllers.NewShiftController(testCase.shiftStore, testCase.soldierStore, testCase.authMiddleware)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, controller)
	}
}

func TestShiftController_NewShiftController__success(t *testing.T) {
	// Act
	controller, err := controllers.NewShiftController(&mocks.MockIShiftStore{}, &mocks.MockISoldierStore{},
		test_utils.AlwaysAllowedJWTMiddleware)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, controller)
}

func TestShiftController_CreateShift__invalid_request_body(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore := &mocks.MockIShiftStore{}
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewShiftController(shiftStore, soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPost, controllers.CreateShiftRoute, test_utils.WrapStructWithReader(t, "invalid"))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestShiftController_CreateShift__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore := &mocks.MockIShiftStore{}
	shiftStore.On("CreateNewShift", mock.MatchedBy(func(arg models.Shift) bool {
		return arg.Name == testShiftName
	})).Return(nil)
	soldierStore := &mocks.MockISoldierStore{}
	soldierStore.On("FindSoldierByID", commanderID).Return([]models.Soldier{testCommander}, nil)
	controller, err := controllers.NewShiftController(shiftStore, soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPost, controllers.CreateShiftRoute, test_utils.WrapStructWithReader(t, testShiftModel))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestShiftController_GetShift__not_found(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore := &mocks.MockIShiftStore{}
	shiftStore.On("FindShiftByID", mock.Anything).Return([]models.Shift{}, nil)
	controller, err := controllers.NewShiftController(shiftStore, &mocks.MockISoldierStore{}, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/shifts/%s", shiftID), nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestShiftController_GetShift__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore := &mocks.MockIShiftStore{}
	shiftStore.On("FindShiftByID", shiftID).Return([]models.Shift{testShiftModel}, nil)
	controller, err := controllers.NewShiftController(shiftStore, &mocks.MockISoldierStore{}, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/shifts/%s", shiftID), nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var respShift models.Shift
	err = json.NewDecoder(resp.Body).Decode(&respShift)
	assert.NoError(t, err)
	assert.Equal(t, testShiftModel, respShift)
}

func TestShiftController_UpdateShift__invalid_request_body(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore := &mocks.MockIShiftStore{}
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewShiftController(shiftStore, soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/shifts/%s", shiftID),
		test_utils.WrapStructWithReader(t, "invalid"))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestShiftController_UpdateShift__not_found(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore := &mocks.MockIShiftStore{}
	shiftStore.On("FindShiftByID", shiftID).Return([]models.Shift{}, nil)
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewShiftController(shiftStore, soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/shifts/%s", shiftID),
		test_utils.WrapStructWithReader(t, testShiftModel))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestShiftController_UpdateShift__success(t *testing.T) {
	// Arrange
	updatedShift := testShiftModel
	updatedShift.Name = "Updated Shift"
	app := fiber.New()
	shiftStoreMock := &mocks.MockIShiftStore{}
	shiftStoreMock.On("FindShiftByID", shiftID).Return([]models.Shift{testShiftModel}, nil)
	shiftStoreMock.On("UpdateShift", mock.MatchedBy(func(arg models.Shift) bool {
		return arg.Name == updatedShift.Name
	})).Return(nil)
	soldierStore := &mocks.MockISoldierStore{}
	soldierStore.On("FindSoldierByID", testShiftModel.Commander.ID).Return([]models.Soldier{testCommander}, nil)
	controller, err := controllers.NewShiftController(shiftStoreMock, soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/shifts/%s", shiftID),
		test_utils.WrapStructWithReader(t, updatedShift))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	shiftStoreMock.AssertExpectations(t)
}

func TestShiftController_DeleteShift__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStoreMock := &mocks.MockIShiftStore{}
	shiftStoreMock.On("DeleteShift", shiftID).Return(nil)
	soldierStore := &mocks.MockISoldierStore{}
	controller, err := controllers.NewShiftController(shiftStoreMock, soldierStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/shifts/%s", shiftID), nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	shiftStoreMock.AssertExpectations(t)
}
