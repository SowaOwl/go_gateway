package model

import "time"

type SystemLog struct {
	ID        uint `gorm:"primary_key"`
	File      string
	Message   string `gorm:"type:text"`
	Stack     string `gorm:"type:text"`
	Line      int
	CreatedAt time.Time
	UpdatedAt time.Time
}
