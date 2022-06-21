package model

import json "github.com/bytedance/sonic"

type GenerateCounter struct {
	ClientID            string `json:"client_id" gorm:"client_id;primaryKey"`
	CreateDate          string `json:"create_date" gorm:"currecreate_datent_date;primaryKey"`
	NumberLinkGenerated int64  `json:"number_link_generated" gorm:"number_link_generated"`
}

func (g GenerateCounter) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(g)
	return bytes, err
}

func (this *GenerateCounter) TableName() string {
	return "shorten_link_generate_counter"
}
