package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	UserID uint   `gorm:"not null" json:"userid"`
	Type   string `gorm:"size:255;not null" json:"type"`
	Fee    uint   `gorm:"not null" json:"fee"`
	Color  string `gorm:"size:255;not null" json:"color"`
}

type User struct {
	gorm.Model
	AuthUserID uint      `gorm:"not null" json:"authuserid"` // AuthUser'ın ID'sini tutacak alan
	Username   string    `gorm:"size:255;not null;unique" json:"username"`
	Products   []Product `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"products"`
}
