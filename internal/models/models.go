package models

import (
	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	Alias string `gorm:"not null;unique;index"`
	Url   string `gorm:"not null"`
}
