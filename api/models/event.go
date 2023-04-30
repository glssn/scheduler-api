package models

import (
	"time"

	"gorm.io/gorm"
)

// Typical event object
type Event struct {
	gorm.Model
	Type              string
	Title             string
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `gorm:"type:TIMESTAMP" json:"end_date"`
	AllDay            bool      `gorm:"default:true" json:"all_day"`
	RecurringType     string    `json:"recurring_type"`
	RecurringInterval uint32    `json:"recurring_interval"`
	User              User
	UserID            int `json:"user_id"`
}

// Typical event metadata object, referring to an Event
type EventMeta struct {
	gorm.Model
	EventID           int
	Event             Event
	RecurringStart    uint64
	RecurringInterval uint32
}
