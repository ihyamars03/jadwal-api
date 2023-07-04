package models

import "time"

type Users struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"unique" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"updatedAt"`
}
