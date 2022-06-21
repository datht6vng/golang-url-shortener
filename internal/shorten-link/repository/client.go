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
func (r *ClientRepository) FindByID(id string) (*model.Client, error) {
	client := new(model.Client)
	err := r.db.Where("client_id = ?", id).First(client).Error
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (r *ClientRepository) Update(client *model.Client) error {
	return r.db.Save(client).Error
}
