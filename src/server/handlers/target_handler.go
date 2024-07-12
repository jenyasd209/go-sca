package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-sca/src/model"
)

const (
	TargetsRoute        = "/targets"
	TargetRoute         = TargetsRoute + "/:id"
	CompleteTargetRoute = TargetRoute + "/complete"
	NotesTargetRoute    = TargetRoute + "/notes"
)

type notesRequest struct {
	NewNotes string `json:"newNotes"`
}

type TargetController interface {
	GetById(uint) (*model.Target, error)
	Complete(uint) error
	UpdateNotes(uint, string) error
	Delete(uint) error
}

type TargetHandler struct {
	targetController TargetController
}

func NewTargetHandler(mc TargetController) *TargetHandler {
	return &TargetHandler{targetController: mc}
}

func (th *TargetHandler) ApplyHandlers(app *fiber.App) {
	app.Get(TargetRoute, th.GetById)
	app.Put(CompleteTargetRoute, th.Complete)
	app.Put(NotesTargetRoute, th.UpdateNotes)
	app.Delete(TargetRoute, th.Delete)
}

func (th *TargetHandler) GetById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid target ID")
	}

	target, err := th.targetController.GetById(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot get target by ID")
	}

	return c.Status(fiber.StatusOK).JSON(target)
}

func (th *TargetHandler) Complete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid target ID")
	}

	err = th.targetController.Complete(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot complete target")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (th *TargetHandler) UpdateNotes(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid target ID")
	}

	notes := &notesRequest{}
	if err := c.BodyParser(&notes); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	err = th.targetController.UpdateNotes(uint(id), notes.NewNotes)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot update target notes")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (th *TargetHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid target ID")
	}

	err = th.targetController.Delete(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot delete target")
	}

	return c.SendStatus(fiber.StatusOK)
}
