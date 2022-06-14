package repository

import (
	"fmt"
	"time"
	"trueid-shorten-link/package/model"

	"gorm.io/gorm"
)

type URLRepository struct {
	db *gorm.DB
}

func (r *URLRepository) Init(db *gorm.DB) *URLRepository {
	r.db = db
	return r
}

func (r *URLRepository) FindByShortURL(shortURL string) (*model.URL, error) {
	url := new(model.URL)
	err := r.db.Where("short_url = ?", shortURL).First(url).Error
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (r *URLRepository) FindByLongURL(longURL string) (*model.URL, error) {
	url := new(model.URL)
	err := r.db.Where("long_url = ?", longURL).First(url).Error
	if err != nil {
		return nil, err
	}
	return url, nil
}

<<<<<<< HEAD
func (this *URLRepository) InsertURL(ID, clientID, shortURL, longURL string, expireTime time.Time) error {
	return this.db.Create(&model.URL{
=======
func (r *URLRepository) InsertURL(ID int64, clientID, shortURL, longURL string, expireTime time.Time) error {
	return r.db.Create(&model.URL{
>>>>>>> master
		ID:         ID,
		ClientID:   clientID,
		ShortURL:   shortURL,
		LongURL:    longURL,
		ExpireTime: expireTime,
	}).Error
<<<<<<< HEAD
	// return this.db.Exec("insert into shorten_link_urls values(?,?,?,?,?)", ID, clientID, shortURL, longURL, expireTime).Error
=======
	// return r.db.Exec("insert into shorten_link_urls values(?,?,?,?,?)", ID, clientID, shortURL, longURL, expireTime).Error
>>>>>>> master
}

func (r *URLRepository) DeleteURL(shortURL, longURL string) error {
	url := new(model.URL)
	if shortURL != "" && longURL != "" {
		return r.db.Where("short_url = ? AND long_url = ?", shortURL, longURL).Delete(&url).Error
	}
	if shortURL == "" && longURL != "" {
		return r.db.Where("long_url = ?", longURL).Delete(&url).Error
	}
	if shortURL != "" && longURL == "" {
		return r.db.Where("short_url = ?", shortURL).Delete(&url).Error
	}
	return r.db.Where("true").Delete(&url).Error
}

func (r *URLRepository) DeleteExpiredURL() error {
	url := new(model.URL)
	return r.db.Where("expire_time < ?", time.Now()).Delete(url).Error
}

<<<<<<< HEAD
func (this *URLRepository) GetMaxID() (string, error) {
	var result string
	if err := this.db.Model(&model.URL{}).Select("max(id)").First(&result).Error; err != nil {
=======
func (r *URLRepository) GetMaxID() (string, error) {
	var result string
	if err := r.db.Model(&model.URL{}).Select("max(id)").First(&result).Error; err != nil {
>>>>>>> master
		fmt.Println(err)
		return "0", err
	}
	return result, nil
}
