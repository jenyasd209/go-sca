package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-sca/src/model"
)

const (
	CatsRoute = "/spy_cats"
	CatRoute  = CatsRoute + "/:id"
)

type salaryRequest struct {
	NewSalary float32 `json:"newSalary"`
}

type CatController interface {
	Create(cat *model.SpyCat) error
	GetById(uint) (*model.SpyCat, error)
	UpdateSalary(uint, float32) error
	Delete(uint) error
	GetAll() ([]*model.SpyCat, error)
}

type CatHandler struct {
	catController CatController
}

func NewCatHandler(cc CatController) *CatHandler {
	return &CatHandler{catController: cc}
}

func (ch *CatHandler) ApplyHandlers(app *fiber.App) {
	app.Post(CatsRoute, ch.Create)
	app.Get(CatsRoute, ch.GetAll)
	app.Get(CatRoute, ch.GetById)
	app.Put(CatRoute, ch.UpdateSalary)
	app.Delete(CatRoute, ch.Delete)
}

func (ch *CatHandler) Create(c *fiber.Ctx) error {
	cat := &model.SpyCat{}

	if err := c.BodyParser(cat); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	err := ch.catController.Create(cat)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot create spy cat")
	}

	return c.Status(fiber.StatusCreated).JSON(cat)
}

func (ch *CatHandler) GetById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid spy cat ID")
	}

	cat, err := ch.catController.GetById(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot get spy cat by ID")
	}

	return c.JSON(cat)
}

func (ch *CatHandler) UpdateSalary(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid spy cat ID")
	}

	salaryReq := &salaryRequest{}
	if err := c.BodyParser(salaryReq); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	err = ch.catController.UpdateSalary(uint(id), salaryReq.NewSalary)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot update salary for spy cat with ID")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (ch *CatHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid spy cat ID")
	}

	err = ch.catController.Delete(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot delete salary for spy cat with ID")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (ch *CatHandler) GetAll(c *fiber.Ctx) error {
	cats, err := ch.catController.GetAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot get all spy cats")
	}

	return c.Status(fiber.StatusCreated).JSON(cats)
}
