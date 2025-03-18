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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDayScheduleController_NewDayScheduleController__error_on_nil_store(t *testing.T) {
	// Act
	controller, err := controllers.NewDayScheduleController(nil, test_utils.AlwaysAllowedJWTMiddleware)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestDayScheduleController_NewDayScheduleController__error_on_nil_auth_middleware(t *testing.T) {
	// Act
	controller, err := controllers.NewDayScheduleController(&mocks.MockIDayStore{}, nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestDayScheduleController_NewDayScheduleController__success(t *testing.T) {
	// Arrange
	dayStore := &mocks.MockIDayStore{}

	// Act
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, controller)
}

func TestDayScheduleController_CreateDaySchedule__invalid_request_body(t *testing.T) {
	// Arrange
	app := fiber.New()
	dayStore := &mocks.MockIDayStore{}
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPost, controllers.CreateDayScheduleRoute, test_utils.WrapStructWithReader(t, "invalid"))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestDayScheduleController_CreateDaySchedule__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	dayStore := &mocks.MockIDayStore{}
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	daySchedule := models.DaySchedule{
		Date: getStrippedUTCDate(),
		Shifts: []models.Shift{
			{ID: "1", Name: "Shift 1"},
		},
	}
	dayStore.On("CreateNewDaySchedule", mock.AnythingOfType("models.DaySchedule")).Return(nil)
	req := httptest.NewRequest(fiber.MethodPost, controllers.CreateDayScheduleRoute, test_utils.WrapStructWithReader(t, daySchedule))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	dayStore.AssertExpectations(t)
}

func TestDayScheduleController_GetDaySchedule__invalid_date_format(t *testing.T) {
	// Arrange
	app := fiber.New()
	dayStore := &mocks.MockIDayStore{}
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodGet, "/day-schedules/invalid-date", nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestDayScheduleController_GetDaySchedule__not_found(t *testing.T) {
	// Arrange
	app := fiber.New()
	dayStore := &mocks.MockIDayStore{}
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	date := getStrippedUTCDate()
	dayStore.On("FindDaySchedule", date).Return([]models.DaySchedule{}, nil)
	req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/day-schedules/%s", date.Format("2006-01-02")), nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	dayStore.AssertExpectations(t)
}

func TestDayScheduleController_GetDaySchedule__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	dayStore := &mocks.MockIDayStore{}
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	date := getStrippedUTCDate()
	daySchedule := models.DaySchedule{
		Date: date,
		Shifts: []models.Shift{
			{ID: "1", Name: "Shift 1"},
		},
	}
	dayStore.On("FindDaySchedule", date).Return([]models.DaySchedule{daySchedule}, nil)
	req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/day-schedules/%s", date.Format("2006-01-02")), nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var respDaySchedule models.DaySchedule
	err = json.NewDecoder(resp.Body).Decode(&respDaySchedule)
	assert.NoError(t, err)
	assert.Equal(t, daySchedule, respDaySchedule)
	dayStore.AssertExpectations(t)
}

func TestDayScheduleController_GetAllDaySchedules__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	dayStore := &mocks.MockIDayStore{}
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	daySchedules := []models.DaySchedule{
		{
			Date: getStrippedUTCDate(),
			Shifts: []models.Shift{
				{ID: "1", Name: "Shift 1"},
			},
		},
		{
			Date: getStrippedUTCDate().Add(24 * time.Hour),
			Shifts: []models.Shift{
				{ID: "2", Name: "Shift 2"},
			},
		},
	}
	dayStore.On("FindAllDaySchedules").Return(daySchedules, nil)
	req := httptest.NewRequest(fiber.MethodGet, controllers.GetAllDaySchedulesRoute, nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var respDaySchedules []models.DaySchedule
	err = json.NewDecoder(resp.Body).Decode(&respDaySchedules)
	assert.NoError(t, err)
	assert.Equal(t, daySchedules, respDaySchedules)
	dayStore.AssertExpectations(t)
}

func TestDayScheduleController_UpdateDaySchedule__invalid_request_body(t *testing.T) {
	// Arrange
	app := fiber.New()
	dayStore := &mocks.MockIDayStore{}
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	date := getStrippedUTCDate()
	req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/day-schedules/%s", date.Format("2006-01-02")), test_utils.WrapStructWithReader(t, "invalid"))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestDayScheduleController_UpdateDaySchedule__invalid_date_format(t *testing.T) {
	// Arrange
	app := fiber.New()
	dayStore := &mocks.MockIDayStore{}
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodPut, "/day-schedules/invalid-date", test_utils.WrapStructWithReader(t, models.DaySchedule{}))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestDayScheduleController_UpdateDaySchedule__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	dayStore := &mocks.MockIDayStore{}
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	date := getStrippedUTCDate()
	daySchedule := models.DaySchedule{
		Date: date,
		Shifts: []models.Shift{
			{ID: "1", Name: "Updated Shift"},
		},
	}
	dayStore.On("UpdateDaySchedule", mock.AnythingOfType("models.DaySchedule")).Return(nil)
	req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/day-schedules/%s", date.Format("2006-01-02")), test_utils.WrapStructWithReader(t, daySchedule))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	dayStore.AssertExpectations(t)
}

func TestDayScheduleController_DeleteDaySchedule__invalid_date_format(t *testing.T) {
	// Arrange
	app := fiber.New()
	dayStore := &mocks.MockIDayStore{}
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	req := httptest.NewRequest(fiber.MethodDelete, "/day-schedules/invalid-date", nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestDayScheduleController_DeleteDaySchedule__success(t *testing.T) {
	// Arrange
	app := fiber.New()
	dayStore := &mocks.MockIDayStore{}
	controller, err := controllers.NewDayScheduleController(dayStore, test_utils.AlwaysAllowedJWTMiddleware)
	require.NoError(t, err)
	err = controller.RegisterRoutes(app)
	require.NoError(t, err)
	date := getStrippedUTCDate()
	dayStore.On("DeleteDaySchedule", date).Return(nil)
	req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/day-schedules/%s", date.Format("2006-01-02")), nil)

	// Act
	resp, err := app.Test(req, test_utils.TestTimeout)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	dayStore.AssertExpectations(t)
}

func getStrippedUTCDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}
