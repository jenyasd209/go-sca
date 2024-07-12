package controller

import (
	"errors"

	"go-sca/src/model"
)

var ErrUnsupportedBreed = errors.New("unsupported breed")

type Repository[T any] interface {
	Create(*T) error
	Get(uint) (*T, error)
	Update(uint, *T, bool) error
	Delete(uint) error
	GetAll() ([]*T, error)
	RawSql(to interface{}, sql string, values ...interface{}) error
	Exec(sql string, values ...interface{}) error
}

type BreedValidator interface {
	Init() error
	Validate(string) bool
}

type SpyCatController struct {
	repo Repository[model.SpyCat]
	bv   BreedValidator
}

func NewSpyCatController(r Repository[model.SpyCat], validator BreedValidator) *SpyCatController {
	return &SpyCatController{
		repo: r,
		bv:   validator,
	}
}

func (scc *SpyCatController) Create(cat *model.SpyCat) error {
	if !scc.bv.Validate(cat.Breed) {
		return ErrUnsupportedBreed
	}

	return scc.repo.Create(cat)
}

func (scc *SpyCatController) GetById(id uint) (*model.SpyCat, error) {
	return scc.repo.Get(id)
}

func (scc *SpyCatController) GetByMissionId(id uint) (*model.SpyCat, error) {
	return scc.repo.Get(id)
}

func (scc *SpyCatController) UpdateSalary(id uint, newSalary float32) error {
	cat, err := scc.repo.Get(id)
	if err != nil {
		return err
	}

	cat.Salary = newSalary
	return scc.repo.Update(id, cat, true)
}

func (scc *SpyCatController) Delete(id uint) error {
	return scc.repo.Delete(id)
}

func (scc *SpyCatController) GetAll() ([]*model.SpyCat, error) {
	return scc.repo.GetAll()
}
