package model

import (
	"gorm.io/gorm"
)

type Mission struct {
	gorm.Model
	ExecutorID uint      `json:"executor_id"`
	Executor   *SpyCat   `json:"executor" gorm:"foreignKey:ExecutorID"`
	Targets    []*Target `json:"targets" gorm:"constraint:OnDelete:CASCADE;"`
	Completed  bool      `json:"completed"`
}
