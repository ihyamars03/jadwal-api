package models

import (
	"time"
)

type Schedule struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	UserId    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Day       string    `json:"day"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdateAt  time.Time `gorm:"autoCreateTime" json:"updatedAt"`
}
