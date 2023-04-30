package models

import "gorm.io/gorm"

// Typical user model
type User struct {
	gorm.Model
	Username string
	Role     string
}
