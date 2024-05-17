package controllers

import (
	"brothers_in_batash/internal/pkg/logging"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type ShiftController struct {
	shiftStore store.IShiftStore
}

func NewShiftController(shiftStore store.IShiftStore) (*ShiftController, error) {
	if shiftStore == nil {
		return nil, errors.New("shiftStore is nil")
	}
	return &ShiftController{shiftStore: shiftStore}, nil
}

func (c *ShiftController) RegisterRoutes(router fiber.Router) error {
	router.Post(CreateShiftRoute, c.createShift)
	router.Get(GetShiftRoute, c.getShift)
	router.Get(GetAllShiftsRoute, c.getAllShifts)
	router.Put(UpdateShiftRoute, c.updateShift)
	router.Delete(DeleteShiftRoute, c.deleteShift)
	return nil
}

func (c *ShiftController) createShift(ctx *fiber.Ctx) error {
	shift := models.Shift{}
	if err := ctx.BodyParser(&shift); err != nil {
		logging.Debug("Could not parse shift creation request body", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := c.shiftStore.CreateNewShift(shift); err != nil {
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
	shifts, err := c.shiftStore.FindAllShifts()
	if err != nil {
		logging.Warning(err, "error on fetching all shifts", nil)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.JSON(shifts)
}

func (c *ShiftController) updateShift(ctx *fiber.Ctx) error {
	shiftID := ctx.Params("id")
	shift := models.Shift{}
	if err := ctx.BodyParser(&shift); err != nil {
		logging.Info("Could not parse shift update request body", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	shift.ID = shiftID
	if err := c.shiftStore.UpdateShift(shift); err != nil {
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
