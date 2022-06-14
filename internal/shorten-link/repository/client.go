package repository

import (
	"trueid-shorten-link/package/model"

	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func (this *ClientRepository) Init(DB *gorm.DB) *ClientRepository {
	this.db = DB
	return this
}
func (this *ClientRepository) FindByAPIKey(apiKey string) (*model.Client, error) {
	client := new(model.Client)
	err := this.db.Where("api_key = ?", apiKey).First(client).Error
	if err != nil {
		return nil, err
	}
	return client, nil

}
