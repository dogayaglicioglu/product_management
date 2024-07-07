package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	UserID uint   `gorm:"size:255;not null;" json:"userid"`
	Type   string `gorm:"size:255;not null;" json:"type"`
	Fee    uint   `gorm:"size:255;not null;" json:"fee"`
	Color  string `gorm:"size:255;not null;" json:"color"`
}
