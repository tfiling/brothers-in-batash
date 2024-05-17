package controllers

import (
	"brothers_in_batash/internal/pkg/logging"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DayScheduleController struct {
	dayStore store.IDayStore
}

func NewDayScheduleController(dayStore store.IDayStore) (*DayScheduleController, error) {
	if dayStore == nil {
		return nil, errors.New("dayStore is nil")
	}
	return &DayScheduleController{dayStore: dayStore}, nil
}

func (c *DayScheduleController) RegisterRoutes(router fiber.Router) error {
	router.Post(CreateDayScheduleRoute, c.createDaySchedule)
	router.Get(GetDayScheduleRoute, c.getDaySchedule)
	router.Get(GetAllDaySchedulesRoute, c.getAllDaySchedules)
	router.Put(UpdateDayScheduleRoute, c.updateDaySchedule)
	router.Delete(DeleteDayScheduleRoute, c.deleteDaySchedule)
	return nil
}

func (c *DayScheduleController) createDaySchedule(ctx *fiber.Ctx) error {
	daySchedule := models.DaySchedule{}
	if err := ctx.BodyParser(&daySchedule); err != nil {
		logging.Debug("Could not parse day schedule creation request body", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := c.dayStore.CreateNewDaySchedule(daySchedule); err != nil {
		logging.Warning(err, "error on creating new day schedule", nil)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusCreated)
}

func (c *DayScheduleController) getDaySchedule(ctx *fiber.Ctx) error {
	dateStr := ctx.Params("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		logging.Debug("Invalid date format", []logging.LogProp{{"date", dateStr}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	daySchedules, err := c.dayStore.FindDaySchedule(date)
	if err != nil {
		logging.Warning(err, "error on fetching day schedule", []logging.LogProp{{"date", dateStr}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if len(daySchedules) == 0 {
		logging.Trace("could not find day schedule", []logging.LogProp{{"date", dateStr}})
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	return ctx.JSON(daySchedules[0])
}

func (c *DayScheduleController) getAllDaySchedules(ctx *fiber.Ctx) error {
	daySchedules, err := c.dayStore.FindAllDaySchedules()
	if err != nil {
		logging.Warning(err, "error on fetching all day schedules", nil)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.JSON(daySchedules)
}

func (c *DayScheduleController) updateDaySchedule(ctx *fiber.Ctx) error {
	dateStr := ctx.Params("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		logging.Debug("Invalid date format", []logging.LogProp{{"date", dateStr}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	daySchedule := models.DaySchedule{}
	if err := ctx.BodyParser(&daySchedule); err != nil {
		logging.Info("Could not parse day schedule update request body", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	daySchedule.Date = date
	if err := c.dayStore.UpdateDaySchedule(daySchedule); err != nil {
		logging.Warning(err, "error on updating day schedule", []logging.LogProp{{"date", dateStr}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (c *DayScheduleController) deleteDaySchedule(ctx *fiber.Ctx) error {
	dateStr := ctx.Params("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		logging.Debug("Invalid date format", []logging.LogProp{{"date", dateStr}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	date = date.UTC()

	if err := c.dayStore.DeleteDaySchedule(date); err != nil {
		logging.Warning(err, "error on deleting day schedule", []logging.LogProp{{"date", dateStr}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}
