package controllers

import (
	"brothers_in_batash/internal/pkg/logging"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// TODO - Implement a concrete type for API requests and responses bodies

type ShiftController struct {
	shiftStore     store.IShiftStore
	soldierStore   store.ISoldierStore
	authMiddleware fiber.Handler
}

func NewShiftController(shiftStore store.IShiftStore, soldierStore store.ISoldierStore, authMiddleware fiber.Handler) (*ShiftController, error) {
	if shiftStore == nil {
		return nil, errors.New("shiftStore is nil")
	}
	if soldierStore == nil {
		return nil, errors.New("soldierStore is nil")
	}
	if authMiddleware == nil {
		return nil, errors.New("authMiddleware is nil")
	}
	return &ShiftController{shiftStore: shiftStore, soldierStore: soldierStore, authMiddleware: authMiddleware}, nil
}

func (c *ShiftController) RegisterRoutes(router fiber.Router) error {
	router.Post(CreateShiftRoute, c.authMiddleware, c.createShift)
	router.Get(GetShiftRoute, c.authMiddleware, c.getShift)
	router.Get(GetAllShiftsRoute, c.authMiddleware, c.getAllShifts)
	router.Put(UpdateShiftRoute, c.authMiddleware, c.updateShift)
	router.Delete(DeleteShiftRoute, c.authMiddleware, c.deleteShift)
	return nil
}

func (c *ShiftController) createShift(ctx *fiber.Ctx) error {
	shiftModel := models.Shift{}
	if err := ctx.BodyParser(&shiftModel); err != nil {
		errStr := err.Error()
		bodyStr := string(ctx.Body())
		logging.Debug("Could not parse shift creation request body", []logging.LogProp{{"error", errStr}, {"body", bodyStr}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := c.shiftStore.CreateNewShift(shiftModel); err != nil {
		logging.Warning(err, "error on creating new shift", nil)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusCreated)
}

func (c *ShiftController) getShift(ctx *fiber.Ctx) error {
	shiftID := ctx.Params("id")
	shifts, err := c.shiftStore.FindShiftByID(shiftID)
	if err != nil {
		logging.Warning(err, "Could not query for shift", []logging.LogProp{{"shiftID", shiftID}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	} else if len(shifts) == 0 {
		logging.Trace("shift not found", []logging.LogProp{{"shiftID", shiftID}})
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	return ctx.JSON(shifts[0])
}

func (c *ShiftController) getAllShifts(ctx *fiber.Ctx) error {
	dbShifts, err := c.shiftStore.FindAllShifts()
	if err != nil {
		logging.Warning(err, "error on fetching all shifts", nil)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	shifts := make([]models.Shift, 0)
	for _, dbShift := range dbShifts {
		shifts = append(shifts, dbShift)
	}
	return ctx.JSON(shifts)
}

func (c *ShiftController) updateShift(ctx *fiber.Ctx) error {
	shiftID := ctx.Params("id")
	updatedShift := models.Shift{}
	if err := ctx.BodyParser(&updatedShift); err != nil {
		logging.Info("Could not parse shift update request body", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	if updatedShift.ID != shiftID {
		logging.Debug("mismatch between shift ID in body and shift ID in URI",
			[]logging.LogProp{{"body_shift_id", updatedShift.ID}, {"uri_shift_id", shiftID}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	if shifts, err := c.shiftStore.FindShiftByID(shiftID); err != nil {
		logging.Info("Could not query existing shift on update", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	} else if len(shifts) == 0 {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	if err := c.shiftStore.UpdateShift(updatedShift); err != nil {
		logging.Warning(err, "error on updating shift", []logging.LogProp{{"shiftID", shiftID}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (c *ShiftController) deleteShift(ctx *fiber.Ctx) error {
	shiftID := ctx.Params("id")
	if err := c.shiftStore.DeleteShift(shiftID); err != nil {
		logging.Warning(err, "error on deleting shift", []logging.LogProp{{"shiftID", shiftID}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}
