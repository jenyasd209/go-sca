package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-sca/src/model"
)

const (
	MissionsRoute        = "/missions"
	MissionRoute         = MissionsRoute + "/:id"
	CompleteMissionRoute = MissionRoute + "/complete"
	AssignMissionRoute   = MissionRoute + "/assign"
	TargetsMissionRoute  = MissionRoute + "/targets"
)

type MissionController interface {
	Create(*model.Mission) error
	GetById(uint) (*model.Mission, error)
	GetAll() ([]*model.Mission, error)
	AssignCat(uint, *model.SpyCat) error
	AddTarget(uint, *model.Target) error
	Complete(uint) error
	Delete(uint) error
}

type MissionHandler struct {
	missionController MissionController
}

func NewMissionHandler(mc MissionController) *MissionHandler {
	return &MissionHandler{missionController: mc}
}

func (mh *MissionHandler) ApplyHandlers(app *fiber.App) {
	app.Post(MissionsRoute, mh.Create)
	app.Get(MissionsRoute, mh.GetAll)
	app.Get(MissionRoute, mh.GetById)
	app.Put(CompleteMissionRoute, mh.Complete)
	app.Put(AssignMissionRoute, mh.AssignCat)
	app.Put(TargetsMissionRoute, mh.AddTarget)
	app.Delete(MissionRoute, mh.Delete)
}

func (mh *MissionHandler) Create(c *fiber.Ctx) error {
	mission := &model.Mission{}
	if err := c.BodyParser(mission); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	err := mh.missionController.Create(mission)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot create new mission")
	}

	return c.Status(fiber.StatusOK).JSON(mission)
}

func (mh *MissionHandler) GetAll(c *fiber.Ctx) error {
	missions, err := mh.missionController.GetAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot get missions")
	}

	return c.Status(fiber.StatusCreated).JSON(missions)
}

func (mh *MissionHandler) GetById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid mission ID")
	}

	mission, err := mh.missionController.GetById(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot get mission by ID")
	}

	return c.Status(fiber.StatusOK).JSON(mission)
}

func (mh *MissionHandler) Complete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid mission ID")
	}

	err = mh.missionController.Complete(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot complete mission")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (mh *MissionHandler) AssignCat(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid mission ID")
	}

	cat := &model.SpyCat{}
	if err := c.BodyParser(cat); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	err = mh.missionController.AssignCat(uint(id), cat)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot assign cat to mission")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (mh *MissionHandler) AddTarget(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid mission ID")
	}

	target := &model.Target{}
	if err := c.BodyParser(target); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	err = mh.missionController.AddTarget(uint(id), target)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot add target to mission")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (mh *MissionHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid mission ID")
	}

	err = mh.missionController.Delete(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot delete mission")
	}

	return c.SendStatus(fiber.StatusOK)
}
