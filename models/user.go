package models

import "gorm.io/gorm"

// https://gorm.io/docs/models.html
// User represents a user resource (for future CRUD with DB).
type User struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Email string `gorm:"not null;uniqueIndex"` // unique constraint on the email column
	Age   int    `gorm:"default:0" json:"age"`
}

// Primary Key: GORM uses a field named ID as the default primary key for each model.
// User model becomes users table with id as primary key.

// predefined struct named gorm.Model is a part of gorm.Model
// it contains the fields ID, CreatedAt, UpdatedAt, DeletedAt
// ID is the primary key
// CreatedAt is the time when the record is created
// UpdatedAt is the time when the record is updated
// DeletedAt is the time when the record is deleted
