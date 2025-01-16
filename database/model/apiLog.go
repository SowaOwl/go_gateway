package model

import "time"

type ApiLog struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	UserID        uint   `gorm:"not null"`
	RequestMethod string `gorm:"not null"`
	Url           string `gorm:"not null"`
	RequestBody   string `gorm:"type:text"`
	RequestHeader string `gorm:"type:text"`
	Ip            string `gorm:"not null"`
	ResponseBody  string `gorm:"type:text"`
	ResponseCode  uint   `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
