package model

import "time"

type ApiLog struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	UserID        uint   `gorm:"not null"`
	RequestMethod string `gorm:"not null;type:varchar(255)"`
	Url           string `gorm:"not null"`
	RequestBody   string `gorm:"type:longtext"`
	RequestHeader string `gorm:"type:longtext"`
	Ip            string `gorm:"not null;type:varchar(255)"`
	ResponseBody  string `gorm:"type:longtext"`
	ResponseCode  uint   `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
