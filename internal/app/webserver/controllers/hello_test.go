package controllers_test

import (
	"brothers_in_batash/internal/app/webserver/controllers"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_HelloController__ok_response(t *testing.T) {
	//Arrange
	app := fiber.New()
	req := httptest.NewRequest("GET", "/hello", nil)
	controller, err := controllers.NewHelloController(app)
	assert.NoError(t, err)
	assert.NotNil(t, controller)
	err = controller.SetupRoutes()
	assert.NoError(t, err)

	//Act
	resp, err := app.Test(req)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(body))
}
