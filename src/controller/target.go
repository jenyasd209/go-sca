package controller

import (
	"errors"

	"go-sca/src/model"
)

var ErrTargetCompleted = errors.New("target is completed")

type TargetController struct {
	repo Repository[model.Target]
	bv   BreedValidator
}

func NewTargetController(r Repository[model.Target]) *TargetController {
	return &TargetController{
		repo: r,
		bv:   NewBreedValidator(),
	}
}

func (tc *TargetController) GetById(id uint) (*model.Target, error) {
	return tc.repo.Get(id)
}

func (tc *TargetController) Complete(id uint) error {
	target, err := tc.repo.Get(id)
	if err != nil {
		return err
	}

	if target.Completed {
		return nil
	}

	target.Completed = true
	return tc.repo.Update(id, target, true)
}

func (tc *TargetController) UpdateNotes(id uint, notes string) error {
	err := tc.repo.Exec(
		`UPDATE targets
				SET notes = ?
				WHERE id IN (
					SELECT t.id
					FROM targets as t JOIN missions m ON t.mission_id = m.id
					WHERE t.completed = 0 AND m.completed = 0 AND t.id = ?
				);`, notes, id)
	if err != nil {
		return err
	}

	return tc.repo.Update(id, nil, false)
}

func (tc *TargetController) Delete(id uint) error {
	target, err := tc.repo.Get(id)
	if err != nil {
		return err
	}

	if target.Completed {
		return ErrTargetCompleted
	}

	return tc.repo.Delete(id)
}
