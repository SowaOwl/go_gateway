package model

import (
	"gorm.io/gorm"
	"time"
)

type WithoutAuthEndpoint struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Value     string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
