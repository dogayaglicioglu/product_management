package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	AuthUserID uint      `gorm:"not null" json:"authuserid"`
	Username   string    `gorm:"size:255;not null;unique" json:"username"`
	Products   []Product `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"products"`
}
