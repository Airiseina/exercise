package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Account  int16  `gorm:"uniqueIndex;not null" json:"account"`
	Username string `gorm:"type:varchar(50);not null" json:"username"`
	Password int16  `gorm:"not null" json:"password"`
}
