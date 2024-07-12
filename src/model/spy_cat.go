package model

import (
	"gorm.io/gorm"
)

type SpyCat struct {
	gorm.Model
	Name       string  `json:"name"`
	Breed      string  `json:"breed"`
	Experience uint8   `json:"experience"`
	Salary     float32 `json:"salary"`
}
