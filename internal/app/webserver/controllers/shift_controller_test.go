package controllers_test

import (
	"brothers_in_batash/internal/app/webserver/controllers"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"brothers_in_batash/internal/pkg/test_utils"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const shiftID = "123"

var testShift = models.Shift{
	ID:   shiftID,
	Name: "Test Shift",
	Type: models.MotorizedPatrolShiftType,
	StartTime: models.TimeOfDay{
		Hour:   0,
		Minute: 0,
	},
	EndTime: models.TimeOfDay{
		Hour:   1,
		Minute: 0,
	},
	Commander: models.Soldier{
		ID:             "123",
		FirstName:      "Gal",
		LastName:       "Tfilin",
		PersonalNumber: "1234567",
		Position:       models.CommanderPosition,
		Roles: []models.SoldierRole{{
			ID:   "1",
			Name: "Commander",
		}},
	},
	AdditionalSoldiers: nil,
}

func TestShiftController_NewShiftController__error_on_nil_store(t *testing.T) {
	// Arrange
	var shiftStore store.IShiftStore

	// Act
	controller, err := controllers.NewShiftController(shiftStore)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestShiftController_NewShiftController__success(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	// Act
	controller, err := controllers.NewShiftController(shiftStore)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, controller)
}

func TestShiftController_CreateShift__invalid_request_body(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)
	controller, err := controllers.NewShiftController(shiftStore)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPost, controllers.CreateShiftRoute, test_utils.WrapStructWithReader(t, "invalid"))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestShiftController_CreateShift__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)
	controller, err := controllers.NewShiftController(shiftStore)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPost, controllers.CreateShiftRoute, test_utils.WrapStructWithReader(t, testShift))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestShiftController_GetShift__not_found(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)
	controller, err := controllers.NewShiftController(shiftStore)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/shifts/%s", shiftID), nil)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestShiftController_GetShift__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)
	controller, err := controllers.NewShiftController(shiftStore)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	err = shiftStore.CreateNewShift(testShift)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/shifts/%s", shiftID), nil)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var respShift models.Shift
	err = json.NewDecoder(resp.Body).Decode(&respShift)
	assert.NoError(t, err)
	assert.Equal(t, testShift, respShift)
}

func TestShiftController_GetAllShifts__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)
	controller, err := controllers.NewShiftController(shiftStore)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	anotherShift := testShift
	anotherShift.ID = "124"
	shifts := []models.Shift{testShift, anotherShift}
	for _, shift := range shifts {
		err = shiftStore.CreateNewShift(shift)
		require.NoError(t, err)
	}
	req := httptest.NewRequest(fiber.MethodGet, controllers.GetAllShiftsRoute, nil)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var respShifts []models.Shift
	err = json.NewDecoder(resp.Body).Decode(&respShifts)
	assert.NoError(t, err)
	assert.ElementsMatch(t, shifts, respShifts)
}

func TestShiftController_UpdateShift__invalid_request_body(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)
	controller, err := controllers.NewShiftController(shiftStore)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/shifts/%s", shiftID),
		test_utils.WrapStructWithReader(t, "invalid"))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestShiftController_UpdateShift__not_found(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)
	controller, err := controllers.NewShiftController(shiftStore)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/shifts/%s", shiftID),
		test_utils.WrapStructWithReader(t, testShift))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestShiftController_UpdateShift__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)
	controller, err := controllers.NewShiftController(shiftStore)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	err = shiftStore.CreateNewShift(testShift)
	require.NoError(t, err)
	updatedShift := testShift
	updatedShift.Name = "Updated Shift"
	req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/shifts/%s", shiftID),
		test_utils.WrapStructWithReader(t, updatedShift))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	foundShift, err := shiftStore.FindShiftByID(shiftID)
	assert.NoError(t, err)
	assert.Len(t, foundShift, 1)
	assert.Equal(t, updatedShift, foundShift[0])
}

func TestShiftController_DeleteShift__not_found(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)
	controller, err := controllers.NewShiftController(shiftStore)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/shifts/%s", shiftID), nil)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestShiftController_DeleteShift__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)
	controller, err := controllers.NewShiftController(shiftStore)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	err = shiftStore.CreateNewShift(testShift)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/shifts/%s", shiftID), nil)

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	foundShifts, err := shiftStore.FindShiftByID(shiftID)
	assert.NoError(t, err)
	assert.Len(t, foundShifts, 0)
}
