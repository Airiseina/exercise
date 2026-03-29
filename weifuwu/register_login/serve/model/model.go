package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Account  string `gorm:"type:varchar(50);uniqueIndex;not null" json:"account"`
	Username string `gorm:"type:varchar(50);not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
}
