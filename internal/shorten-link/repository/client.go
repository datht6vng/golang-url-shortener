package repository

import (
	"trueid-shorten-link/package/model"

	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func (r *ClientRepository) Init(DB *gorm.DB) *ClientRepository {
	r.db = DB
	return r
}
func (r *ClientRepository) FindByAPIKey(apiKey string) (*model.Client, error) {
	client := new(model.Client)
	err := r.db.Where("api_key = ?", apiKey).First(client).Error
	if err != nil {
		return nil, err
	}
	return client, nil
}
