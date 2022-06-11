package repository

import (
	"time"
	"trueid-shorten-link/package/model"

	"gorm.io/gorm"
)

type URLRepository struct {
	db *gorm.DB
}

func (this *URLRepository) Init(db *gorm.DB) *URLRepository {
	this.db = db
	return this
}

func (this *URLRepository) FindByShortURL(shortURL string) (*model.URL, error) {
	url := new(model.URL)
	err := this.db.Where("short_url = ?", shortURL).First(url).Error
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (this *URLRepository) FindByLongURL(longURL string) (*model.URL, error) {
	url := new(model.URL)
	err := this.db.Where("long_url = ?", longURL).First(url).Error
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (this *URLRepository) InsertURL(ID, clientID, shortURL, longURL string, expireTime time.Time) error {
	return this.db.Exec("insert into urls values(?,?,?,?,?)", ID, clientID, shortURL, longURL, expireTime).Error
}

func (this *URLRepository) DeleteURL(shortURL, longURL string) error {
	url := new(model.URL)
	if shortURL != "" && longURL != "" {
		return this.db.Where("short_url = ? AND long_url = ?", shortURL, longURL).Delete(&url).Error
	}
	if shortURL == "" && longURL != "" {
		return this.db.Where("long_url = ?", longURL).Delete(&url).Error
	}
	if shortURL != "" && longURL == "" {
		return this.db.Where("short_url = ?", shortURL).Delete(&url).Error
	}
	return this.db.Where("true").Delete(&url).Error
}

func (this *URLRepository) DeleteExpiredURL() error {
	url := new(model.URL)
	return this.db.Where("expire_time < ?", time.Now()).Delete(url).Error
}

func (this *URLRepository) GetMaxID() (string, error) {
	var result interface{}
	err := this.db.Table("urls").Select("max(id)").Row().Scan(&result)
	if result == nil {
		return "0", err
	}
	return string(result.([]byte)), err
}
