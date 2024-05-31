package controllers

import (
	"brothers_in_batash/internal/pkg/logging"
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"errors"

	"github.com/gofiber/fiber/v2"
)

// TODO - Implement a concrete type for API requests and responses bodies(currently using the actual models)

type SoldierController struct {
	soldierStore   store.ISoldierStore
	authMiddleware fiber.Handler
}

func NewSoldierController(soldierStore store.ISoldierStore, authMiddleware fiber.Handler) (*SoldierController, error) {
	if soldierStore == nil {
		return nil, errors.New("soldierStore is nil")
	}
	if authMiddleware == nil {
		return nil, errors.New("authMiddleware is nil")
	}
	return &SoldierController{soldierStore: soldierStore, authMiddleware: authMiddleware}, nil
}

func (c *SoldierController) RegisterRoutes(router fiber.Router) error {
	router.Post(CreateSoldierRoute, c.createSoldier)
	router.Get(GetSoldierRoute, c.getSoldier)
	router.Get(GetAllSoldiersRoute, c.getAllSoldiers)
	router.Put(UpdateSoldierRoute, c.updateSoldier)
	router.Delete(DeleteSoldierRoute, c.deleteSoldier)
	return nil
}

func (c *SoldierController) createSoldier(ctx *fiber.Ctx) error {
	soldier := models.Soldier{}
	if err := ctx.BodyParser(&soldier); err != nil {
		logging.Debug("Could not parse soldier creation request body", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := c.soldierStore.CreateNewSoldier(soldier); err != nil {
		logging.Warning(err, "error on creating new soldier", nil)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusCreated)
}

func (c *SoldierController) getSoldier(ctx *fiber.Ctx) error {
	soldierID := ctx.Params("id")
	soldiers, err := c.soldierStore.FindSoldierByID(soldierID)
	if err != nil {
		logging.Warning(err, "could not query for soldier", []logging.LogProp{{"soldierID", soldierID}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	} else if len(soldiers) == 0 {
		logging.Trace("Soldier not found", []logging.LogProp{{"soldierID", soldierID}})
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	return ctx.JSON(soldiers[0])
}

func (c *SoldierController) getAllSoldiers(ctx *fiber.Ctx) error {
	soldiers, err := c.soldierStore.FindAllSoldiers()
	if err != nil {
		logging.Warning(err, "error on fetching all soldiers", nil)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.JSON(soldiers)
}

func (c *SoldierController) updateSoldier(ctx *fiber.Ctx) error {
	soldierID := ctx.Params("id")
	soldier := models.Soldier{}
	if err := ctx.BodyParser(&soldier); err != nil {
		logging.Info("Could not parse soldier update request body", []logging.LogProp{{"error", err.Error()}})
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	soldier.ID = soldierID
	if err := c.soldierStore.UpdateSoldier(soldier); err != nil {
		logging.Warning(err, "error on updating soldier", []logging.LogProp{{"soldierID", soldierID}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (c *SoldierController) deleteSoldier(ctx *fiber.Ctx) error {
	soldierID := ctx.Params("id")
	if err := c.soldierStore.DeleteSoldier(soldierID); err != nil {
		logging.Warning(err, "error on deleting soldier", []logging.LogProp{{"soldierID", soldierID}})
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}
