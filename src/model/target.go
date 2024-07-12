package model

import (
	"gorm.io/gorm"
)

type Target struct {
	gorm.Model
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	Completed bool   `json:"completed"`
	MissionID uint   `json:"mission_id"`
}
