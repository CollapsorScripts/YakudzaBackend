package migrator

import (
	"Yakudza/pkg/database/models"
	"errors"
	"gorm.io/gorm"
)

func Migrations(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.Links{})
	if err != nil {
		return errors.New("Migration error: " + err.Error())
	}

	return nil
}
