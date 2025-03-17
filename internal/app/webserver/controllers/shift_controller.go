package controllers

import (
	"brothers_in_batash/internal/app/webserver/api"
	"brothers_in_batash/internal/pkg/logging"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"brothers_in_batash/internal/pkg/utils"

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
	reqBody := api.Shift{}
	if err := ctx.BodyParser(&reqBody); err != nil {
		errStr := err.Error()
		bodyStr := string(ctx.Body())
		logging.Debug("Could not parse shift creation request body", []logging.LogProp{{"error", errStr}, {"body", bodyStr}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	shiftModel, err := c.convertShiftApiReqToModel(ctx, reqBody)
	if err != nil {
		return err
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
	return ctx.JSON(c.convertShiftModelToApiTYpe(shifts[0]))
}

func (c *ShiftController) getAllShifts(ctx *fiber.Ctx) error {
	dbShifts, err := c.shiftStore.FindAllShifts()
	if err != nil {
		logging.Warning(err, "error on fetching all shifts", nil)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	shifts := make([]api.Shift, 1)
	for _, dbShift := range dbShifts {
		shifts = append(shifts, c.convertShiftModelToApiTYpe(dbShift))
	}
	return ctx.JSON(shifts)
}

func (c *ShiftController) updateShift(ctx *fiber.Ctx) error {
	shiftID := ctx.Params("id")
	apiShift := api.Shift{}
	if err := ctx.BodyParser(&apiShift); err != nil {
		logging.Info("Could not parse shift update request body", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	if apiShift.ID != shiftID {
		logging.Debug("mismatch between shift ID in body and shift ID in URI",
			[]logging.LogProp{{"body_shift_id", apiShift.ID}, {"uri_shift_id", shiftID}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	if shifts, err := c.shiftStore.FindShiftByID(shiftID); err != nil {
		logging.Info("Could not query existing shift on update", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	} else if len(shifts) == 0 {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	shift, err := c.convertShiftApiReqToModel(ctx, apiShift)
	if err != nil {
		logging.Debug("Could not convert API shift to shift model", []logging.LogProp{{"error", err.Error()}})
		return err
	}
	if err = c.shiftStore.UpdateShift(shift); err != nil {
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

func (c *ShiftController) convertShiftApiReqToModel(ctx *fiber.Ctx, shiftApiReq api.Shift) (models.Shift, error) {
	shiftModel := models.Shift{
		ID: utils.NewEntityID(),
		ShiftTime: models.ShiftTime{
			StartTime: models.TimeOfDay{
				Hour:   shiftApiReq.StartTimeHour,
				Minute: shiftApiReq.StartTimeMinute,
			},
			EndTime: models.TimeOfDay{
				Hour:   shiftApiReq.EndTimeHour,
				Minute: shiftApiReq.EndTimeMinute,
			},
		},
		Name:            shiftApiReq.Name,
		Type:            models.ShiftType(shiftApiReq.Type),
		Description:     shiftApiReq.Description,
		ShiftTemplateID: shiftApiReq.ShiftTemplateID,
	}
	return c.populateSoldiersFromIDs(ctx, shiftApiReq, shiftModel)
}

func (c *ShiftController) populateSoldiersFromIDs(ctx *fiber.Ctx, shiftApiReq api.Shift, shiftModel models.Shift) (models.Shift, error) {
	foundSoldiers, err := c.soldierStore.FindSoldierByID(shiftApiReq.CommanderSoldierID)
	if err != nil {
		logging.Warning(err, "could not query shift commander", []logging.LogProp{{"commander_id", shiftApiReq.CommanderSoldierID}})
		return shiftModel, ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if len(foundSoldiers) == 0 {
		logging.Debug("could not query shift commander", []logging.LogProp{{"commander_id", shiftApiReq.CommanderSoldierID}})
		return shiftModel, ctx.SendStatus(fiber.StatusBadRequest)
	}

	shiftModel.Commander = foundSoldiers[0]

	for _, soldierID := range shiftApiReq.AdditionalSoldiersIDs {
		foundSoldiers, err = c.soldierStore.FindSoldierByID(soldierID)
		if err != nil {
			logging.Warning(err, "could not query additional soldier", []logging.LogProp{{"soldier_id", soldierID}})
			return shiftModel, ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if len(foundSoldiers) == 0 {
			logging.Debug("could not query additional soldier", []logging.LogProp{{"soldier_id", soldierID}})
			return shiftModel, ctx.SendStatus(fiber.StatusBadRequest)
		}
		shiftModel.AdditionalSoldiers = append(shiftModel.AdditionalSoldiers, foundSoldiers[0])
	}

	return shiftModel, nil
}

func (c *ShiftController) convertShiftModelToApiTYpe(shift models.Shift) api.Shift {
	var soldierIDs = make([]string, 0)
	for _, soldier := range shift.AdditionalSoldiers {
		soldierIDs = append(soldierIDs, soldier.ID)
	}
	return api.Shift{
		StartTimeHour:         shift.StartTime.Hour,
		StartTimeMinute:       shift.StartTime.Minute,
		EndTimeHour:           shift.EndTime.Hour,
		EndTimeMinute:         shift.EndTime.Minute,
		ID:                    shift.ID,
		Name:                  shift.Name,
		Type:                  api.ShiftType(shift.Type),
		CommanderSoldierID:    shift.Commander.ID,
		AdditionalSoldiersIDs: soldierIDs,
		Description:           shift.Description,
		ShiftTemplateID:       shift.ShiftTemplateID,
	}
}
