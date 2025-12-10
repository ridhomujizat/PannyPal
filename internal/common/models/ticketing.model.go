package models

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	Title       string `gorm:"type:varchar(200);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Status      string `gorm:"type:varchar(50);not null" json:"status"`
}
