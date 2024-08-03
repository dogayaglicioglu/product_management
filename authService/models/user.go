package models

import "gorm.io/gorm"

type AuthUser struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

type RequestPayload struct {
	NewUsername string `json:"new_username"`
	OldUsername string `json:"old_username"`
}
