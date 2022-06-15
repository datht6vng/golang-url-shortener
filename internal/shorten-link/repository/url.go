package repository

import (
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

func (r *URLRepository) InsertURL(ID int64, clientID, shortURL, longURL string, expireTime time.Time) error {
	return r.db.Create(&model.URL{
		ID:         ID,
		ClientID:   clientID,
		ShortURL:   shortURL,
		LongURL:    longURL,
		ExpireTime: expireTime,
	}).Error
	// return r.db.Exec("insert into shorten_link_urls values(?,?,?,?,?)", ID, clientID, shortURL, longURL, expireTime).Error
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

func (r *URLRepository) GetMaxID() (string, error) {
	var result string
	if err := r.db.Model(&model.URL{}).Select("max(id)").First(&result).Error; err != nil {
		return "0", err
	}
	return result, nil
}
func (r *URLRepository) CountLinkGenerated(clientID string) (int64, error) {
	var result int64
	if err := r.db.Model(&model.URL{}).Select("count(*)").Where(`
		client_id = ? and cast(expire_time as date) = ? 
	`, clientID, time.Now().AddDate(0, 0, 3).Format("2006-01-02")).First(&result).Error; err != nil {
		return 0, err
	}
	return result, nil
}
