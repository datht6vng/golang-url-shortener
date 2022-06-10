package repository

import (
	"time"
	"trueid-shorten-link/package/model"

	"gorm.io/gorm"
)

type UrlRepository struct {
	DB *gorm.DB
}

func (this *UrlRepository) Init(DB *gorm.DB) *UrlRepository {
	this.DB = DB
	return this
}

func (this *UrlRepository) FindShortUrl(shortUrl string) (*model.Url, error) {
	url := new(model.Url)
	err := this.DB.Where("short_url = ?", shortUrl).First(&url).Error
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (this *UrlRepository) FindLongUrl(longUrl string) (*model.Url, error) {
	url := new(model.Url)
	err := this.DB.Where("long_url = ?", longUrl).First(&url).Error
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (this *UrlRepository) InsertUrl(ID, user, shortUrl, longUrl string, expireTime time.Time) error {
	return this.DB.Exec("insert into urls values(?,?,?,?,?)", ID, user, shortUrl, longUrl, expireTime).Error
}

func (this *UrlRepository) DeleteUrl(shortUrl, longUrl string) error {
	url := new(model.Url)
	if shortUrl != "" && longUrl != "" {
		return this.DB.Where("short_url = ? AND long_url = ?", shortUrl, longUrl).Delete(&url).Error
	}
	if shortUrl == "" && longUrl != "" {
		return this.DB.Where("long_url = ?", longUrl).Delete(&url).Error
	}
	if shortUrl != "" && longUrl == "" {
		return this.DB.Where("short_url = ?", shortUrl).Delete(&url).Error
	}
	return this.DB.Where("true").Delete(&url).Error
}

func (this *UrlRepository) DeleteExpiredUrl() error {
	url := new(model.Url)
	return this.DB.Where("expire_time < ?", time.Now()).Delete(url).Error
}

func (this *UrlRepository) GetMaxID() (string, error) {
	var result interface{}
	err := this.DB.Table("urls").Select("max(id)").Row().Scan(&result)
	if result == nil {
		return "0", err
	}
	return string(result.([]byte)), err
}
