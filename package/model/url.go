package model

import "time"

type Url struct {
	ID         int64     `json:"ID" xml:"ID" form:"ID" gorm:"column=id:PRIMARY;index"`
	User       string    `json:"user" xml:"user" form:"user" gorm:"user"`
	ShortUrl   string    `json:"shortUrl" xml:"shortUrl" form:"shortUrl" gorm:"short_url;index"`
	LongUrl    string    `json:"longUrl"  xml:"longUrl"  form:"longUrl" gorm:"long_url;index"`
	ExpireTime time.Time `json:"expireTime" xml:"expireTime" form:"expireTime" gorm:"expire_time"`
}
