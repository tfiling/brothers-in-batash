package controllers_test

import (
	"brothers_in_batash/internal/app/webserver/controllers"
	"brothers_in_batash/internal/pkg/mocks"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/test_utils"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const shiftTemplateID = "123"

var testShiftTemplate = models.ShiftTemplate{
	ID:          shiftTemplateID,
	Name:        "Test Shift Template",
	Description: "Test Description",
	PersonnelRequirement: models.PersonnelRequirement{
		SoldierRoleToCount: map[string]int{
			"Commander": 1,
		},
	},
	DaysOfOccurrences: map[time.Weekday][]models.ShiftTime{
		time.Monday: {
			{
				StartTime: models.TimeOfDay{Hour: 0, Minute: 0},
				EndTime:   models.TimeOfDay{Hour: 1, Minute: 0},
			},
		},
	},
}

func TestShiftTemplateController_NewShiftTemplateController__error_on_nil_store(t *testing.T) {
	// Act
	controller, err := controllers.NewShiftTemplateController(nil, test_utils.AlwaysAllowedJWTMiddleware)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestShiftTemplateController_NewShiftTemplateController__error_on_nil_middleware(t *testing.T) {
	// Arrange
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)

	// Act
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestShiftTemplateController_NewShiftTemplateController__success(t *testing.T) {
	// Arrange
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)

	// Act
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, test_utils.AlwaysAllowedJWTMiddleware)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, controller)
}

func TestShiftTemplateController_GetShiftTemplate__not_found(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)

	shiftTemplateStore.On("FindShiftTemplateByID", shiftTemplateID).Return([]models.ShiftTemplate{}, nil)

	req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/shift-templates/%s", shiftTemplateID), nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	shiftTemplateStore.AssertExpectations(t)
}

func TestShiftTemplateController_GetShiftTemplate__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)

	shiftTemplateStore.On("FindShiftTemplateByID", shiftTemplateID).Return([]models.ShiftTemplate{testShiftTemplate}, nil)

	req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/shift-templates/%s", shiftTemplateID), nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var respShiftTemplate models.ShiftTemplate
	err = json.NewDecoder(resp.Body).Decode(&respShiftTemplate)
	assert.NoError(t, err)
	assert.Equal(t, testShiftTemplate, respShiftTemplate)
	shiftTemplateStore.AssertExpectations(t)
}

func TestShiftTemplateController_GetAllShiftTemplates__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)

	anotherShiftTemplate := testShiftTemplate
	anotherShiftTemplate.ID = "124"
	shiftTemplates := []models.ShiftTemplate{testShiftTemplate, anotherShiftTemplate}

	shiftTemplateStore.On("FindAllShiftsTemplate").Return(shiftTemplates, nil)

	req := httptest.NewRequest(fiber.MethodGet, controllers.GetAllShiftTemplatesRoute, nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var respShiftTemplates []models.ShiftTemplate
	err = json.NewDecoder(resp.Body).Decode(&respShiftTemplates)
	assert.NoError(t, err)
	assert.ElementsMatch(t, shiftTemplates, respShiftTemplates)
	shiftTemplateStore.AssertExpectations(t)
}

func TestShiftTemplateController_DeleteShiftTemplate__not_found(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)

	shiftTemplateStore.On("DeleteShiftTemplate", shiftTemplateID).Return(fmt.Errorf("shift template not found"))

	req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/shift-templates/%s", shiftTemplateID), nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	shiftTemplateStore.AssertExpectations(t)
}

func TestShiftTemplateController_DeleteShiftTemplate__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)

	shiftTemplateStore.On("DeleteShiftTemplate", shiftTemplateID).Return(nil)

	req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/shift-templates/%s", shiftTemplateID), nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	shiftTemplateStore.AssertExpectations(t)
}

func TestShiftTemplateController_CreateShiftTemplate__invalid_request_body(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)

	req := httptest.NewRequest(fiber.MethodPost, controllers.CreateShiftTemplateRoute, test_utils.WrapStructWithReader(t, "invalid"))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestShiftTemplateController_CreateShiftTemplate__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)

	shiftTemplateStore.On("CreateNewShiftTemplate", testShiftTemplate).Return(nil)

	req := httptest.NewRequest(fiber.MethodPost, controllers.CreateShiftTemplateRoute, test_utils.WrapStructWithReader(t, testShiftTemplate))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	shiftTemplateStore.AssertExpectations(t)
}

func TestShiftTemplateController_UpdateShiftTemplate__invalid_request_body(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)

	req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/shift-templates/%s", shiftTemplateID),
		test_utils.WrapStructWithReader(t, "invalid"))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestShiftTemplateController_UpdateShiftTemplate__not_found(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)

	shiftTemplateStore.On("UpdateShiftTemplate", testShiftTemplate).Return(fmt.Errorf("shift template not found"))

	req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/shift-templates/%s", shiftTemplateID),
		test_utils.WrapStructWithReader(t, testShiftTemplate))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	shiftTemplateStore.AssertExpectations(t)
}

func TestShiftTemplateController_UpdateShiftTemplate__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	shiftTemplateStore := new(mocks.MockIShiftTemplateStore)
	controller, err := controllers.NewShiftTemplateController(shiftTemplateStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)

	updatedShiftTemplate := testShiftTemplate
	updatedShiftTemplate.Name = "Updated Shift Template"

	shiftTemplateStore.On("UpdateShiftTemplate", updatedShiftTemplate).Return(nil)

	req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/shift-templates/%s", shiftTemplateID),
		test_utils.WrapStructWithReader(t, updatedShiftTemplate))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	shiftTemplateStore.AssertExpectations(t)
}
