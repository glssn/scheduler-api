package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Type              string
	Title             string
	StartDate         time.Time
	EndDate           time.Time `gorm:"type:TIMESTAMP"`
	AllDay            bool      `gorm:"default:true"`
	RecurringType     string
	RecurringInterval uint32
	User              User
	UserID            int
}

type EventMeta struct {
	gorm.Model
	EventID           int
	Event             Event
	RecurringStart    uint64
	RecurringInterval uint32
}
