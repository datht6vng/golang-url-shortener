package migration

import (
	"trueid-shorten-link/package/model"

	"gorm.io/gorm"
)

func CreateTable(db *gorm.DB) {
	db.AutoMigrate(&model.Url{})
}
