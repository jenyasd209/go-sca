package database

import (
	"go-sca/src/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(name string, cfg *gorm.Config) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(name), cfg)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Target{}, &model.Mission{}, &model.SpyCat{})
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}
