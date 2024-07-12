package controller

import (
	"errors"

	"go-sca/src/model"
)

var (
	ErrCatAssigned      = errors.New("a cat is assigned")
	ErrMissionCompleted = errors.New("mission completed")
)

type MissionController struct {
	missionRepo Repository[model.Mission]
	bv          BreedValidator
}

func NewMissionController(mr Repository[model.Mission]) *MissionController {
	return &MissionController{
		missionRepo: mr,
		bv:          NewBreedValidator(),
	}
}

func (mc *MissionController) Create(mission *model.Mission) error {
	return mc.missionRepo.Create(mission)
}

func (mc *MissionController) GetById(id uint) (*model.Mission, error) {
	mission, err := mc.missionRepo.Get(id)
	if err != nil {
		return nil, err
	}

	err = mc.readMissionRelatedData(mission)
	if err != nil {
		return nil, err
	}

	return mission, nil
}

func (mc *MissionController) GetAll() ([]*model.Mission, error) {
	missions, err := mc.missionRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, mission := range missions {
		err = mc.readMissionRelatedData(mission)
		if err != nil {
			return nil, err
		}
	}

	return missions, nil
}

func (mc *MissionController) AssignCat(id uint, cat *model.SpyCat) error {
	mission, err := mc.missionRepo.Get(id)
	if err != nil {
		return err
	}

	mission.Executor = cat
	return mc.missionRepo.Update(id, mission, true)
}

func (mc *MissionController) AddTarget(id uint, t *model.Target) error {
	mission, err := mc.GetById(id)
	if err != nil {
		return err
	}

	if mission.Completed {
		return ErrMissionCompleted
	}

	t.MissionID = mission.ID
	return mc.missionRepo.Update(id, mission, true)
}

func (mc *MissionController) Complete(id uint) error {
	mission, err := mc.missionRepo.Get(id)
	if err != nil {
		return err
	}

	if mission.Completed {
		return nil
	}

	//TODO not sure that targets should be completed too
	//for _, target := range mission.Targets {
	//	target.Completed = true
	//}

	mission.Completed = true
	return mc.missionRepo.Update(id, mission, true)
}

func (mc *MissionController) Delete(id uint) error {
	m, err := mc.GetById(id)
	if err != nil {
		return err
	}

	if m.Executor != nil {
		return ErrCatAssigned
	}

	return mc.missionRepo.Delete(id)
}

func (mc *MissionController) readMissionRelatedData(mission *model.Mission) error {
	if mission.ExecutorID != 0 {
		executor := &model.SpyCat{}
		err := mc.missionRepo.RawSql(executor, "SELECT * FROM spy_cats WHERE id = ?", mission.ExecutorID)
		if err != nil {
			return err
		}

		mission.Executor = executor
	}

	var targets []*model.Target
	err := mc.missionRepo.RawSql(&targets, "SELECT * FROM targets WHERE mission_id = ?", mission.ID)
	if err != nil {
		return err
	}
	mission.Targets = targets

	return nil
}
