package model

import "time"

type URL struct {
	ID         string    `json:"id" xml:"id" form:"id" gorm:"id;primaryKey;index"`
	ClientID   string    `json:"client_id" xml:"client_id" form:"client_id" gorm:"client_id"`
	ShortURL   string    `json:"shortURL" xml:"shortURL" form:"shortURL" gorm:"short_url;index"`
	LongURL    string    `json:"longURL"  xml:"longURL"  form:"longURL" gorm:"long_url;index"`
	ExpireTime time.Time `json:"expireTime" xml:"expireTime" form:"expireTime" gorm:"expire_time"`
}

func (m *URL) TableName() string {
	return `shorten_link_urls`
}
