package repository

import (
	"trueid-shorten-link/package/model"

	"gorm.io/gorm"
)

type GenerateCounterRepository struct {
	db *gorm.DB
}

func (this *GenerateCounterRepository) Init(DB *gorm.DB) *GenerateCounterRepository {
	this.db = DB
	return this
}
func (this *GenerateCounterRepository) Insert(client_id, createDate string, numberLinkGenerated int64) error {
	return this.db.Create(&model.GenerateCounter{
		ClientID:            client_id,
		CreateDate:          createDate,
		NumberLinkGenerated: numberLinkGenerated,
	}).Error
}

func (this *GenerateCounterRepository) Update(client_id, createDate string, numberLinkGenerated int64) error {
	return this.db.Model(&model.GenerateCounter{
		ClientID:            client_id,
		CreateDate:          createDate,
		NumberLinkGenerated: numberLinkGenerated,
	}).UpdateColumn("number_link_generated", gorm.Expr("number_link_generated + ?", numberLinkGenerated)).Error
}
