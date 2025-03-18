package controllers

import (
	"brothers_in_batash/internal/pkg/logging"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type ShiftTemplateController struct {
	shiftTemplateStore store.IShiftTemplateStore
	authMiddleware     fiber.Handler
}

func NewShiftTemplateController(shiftTemplateStore store.IShiftTemplateStore, authMiddleware fiber.Handler) (*ShiftTemplateController, error) {
	if shiftTemplateStore == nil {
		return nil, errors.New("shiftTemplateStore is nil")
	}
	if authMiddleware == nil {
		return nil, errors.New("authMiddleware is nil")
	}
	return &ShiftTemplateController{shiftTemplateStore: shiftTemplateStore, authMiddleware: authMiddleware}, nil
}

func (c *ShiftTemplateController) RegisterRoutes(router fiber.Router) error {
	router.Post(CreateShiftTemplateRoute, c.authMiddleware, c.createShiftTemplate)
	router.Get(GetShiftTemplateRoute, c.authMiddleware, c.getShiftTemplate)
	router.Get(GetAllShiftTemplatesRoute, c.authMiddleware, c.getAllShiftTemplates)
	router.Put(UpdateShiftTemplateRoute, c.authMiddleware, c.updateShiftTemplate)
	router.Delete(DeleteShiftTemplateRoute, c.authMiddleware, c.deleteShiftTemplate)
	return nil
}

func (c *ShiftTemplateController) createShiftTemplate(ctx *fiber.Ctx) error {
	shiftTemplate := models.ShiftTemplate{}
	if err := ctx.BodyParser(&shiftTemplate); err != nil {
		logging.Debug("Could not parse shift template creation request body", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := c.shiftTemplateStore.CreateNewShiftTemplate(shiftTemplate); err != nil {
		logging.Warning(err, "error on creating new shift template", nil)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusCreated)
}

func (c *ShiftTemplateController) getShiftTemplate(ctx *fiber.Ctx) error {
	shiftTemplateID := ctx.Params("id")
	shiftTemplates, err := c.shiftTemplateStore.FindShiftTemplateByID(shiftTemplateID)
	if err != nil {
		logging.Warning(err, "Could not query for shift template", []logging.LogProp{{"shiftTemplateID", shiftTemplateID}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	} else if len(shiftTemplates) == 0 {
		logging.Trace("shift template not found", []logging.LogProp{{"shiftTemplateID", shiftTemplateID}})
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	return ctx.JSON(shiftTemplates[0])
}

func (c *ShiftTemplateController) getAllShiftTemplates(ctx *fiber.Ctx) error {
	shiftTemplates, err := c.shiftTemplateStore.FindAllShiftsTemplate()
	if err != nil {
		logging.Warning(err, "error on fetching all shift templates", nil)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.JSON(shiftTemplates)
}

func (c *ShiftTemplateController) updateShiftTemplate(ctx *fiber.Ctx) error {
	shiftTemplateID := ctx.Params("id")
	shiftTemplate := models.ShiftTemplate{}
	if err := ctx.BodyParser(&shiftTemplate); err != nil {
		logging.Info("Could not parse shift template update request body", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	shiftTemplate.ID = shiftTemplateID
	if err := c.shiftTemplateStore.UpdateShiftTemplate(shiftTemplate); err != nil {
		logging.Warning(err, "error on updating shift template", []logging.LogProp{{"shiftTemplateID", shiftTemplateID}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (c *ShiftTemplateController) deleteShiftTemplate(ctx *fiber.Ctx) error {
	shiftTemplateID := ctx.Params("id")
	if err := c.shiftTemplateStore.DeleteShiftTemplate(shiftTemplateID); err != nil {
		logging.Warning(err, "error on deleting shift template", []logging.LogProp{{"shiftTemplateID", shiftTemplateID}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}
