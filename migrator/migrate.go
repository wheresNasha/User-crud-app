package migrator

// https://gorm.io/docs/migration.html

import (
	"go-crud/initializers"
	"go-crud/models"
)

func Migrate() {
	initializers.DB.AutoMigrate(&models.User{})
	// AutoMigrate is a method on the DB type struct,
	// so AutoMigrate is a method on the initializers.DB variable
	//func (db *DB) AutoMigrate(dst ...interface{}) error -> takes db Receiver
	// https://pkg.go.dev/gorm.io/gorm@v1.25.10#DB.AutoMigrate
	// AutoMigrate run auto migration for given models
	// it will create tables, missing foreign keys, and create missing indexes
}
