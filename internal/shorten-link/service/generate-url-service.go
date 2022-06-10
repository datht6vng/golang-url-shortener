package service

import (
	"fmt"
	"time"
	"trueid-shorten-link/internal/shorten-link/repository"
	"trueid-shorten-link/package/encryption"
	"trueid-shorten-link/package/metrics"
	_redis "trueid-shorten-link/package/redis"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

type UrlData struct {
	User string `json:"user" xml:"user" form:"user"`
	Url  string `json:"url" xml:"url" form:"url"`
}

type GenerateUrlService struct {
	urlRepository *repository.UrlRepository
	redis         *_redis.Redis
	metrics       *metrics.Metrics
}

func (this *GenerateUrlService) Init(urlRepository *repository.UrlRepository, redis *_redis.Redis, metrics *metrics.Metrics) *GenerateUrlService {
	this.urlRepository = urlRepository
	this.redis = redis
	this.metrics = metrics
	return this
}
func (this *GenerateUrlService) GetNextID() string {
	if currentID := this.redis.Get("CurrentID"); currentID == "" {
		pipe := this.redis.TxPipeline()
		resultID := ""
		this.redis.Watch(func(tx *redis.Tx) error {
			maxID, _ := this.urlRepository.GetMaxID()
			pipe.Set("CurrentID", maxID, -1)
			cmd := pipe.Incr("CurrentID")
			if _, err := pipe.Exec(); err != nil && err != redis.Nil {
				return err
			}
			resultID = fmt.Sprint(cmd.Val())
			return nil
		}, "CurrentID")
		return resultID
	}
	return fmt.Sprint(this.redis.Incr("CurrentID"))
}
func (this *GenerateUrlService) GenerateUrl(inputData *UrlData) (string, error) {
	go this.metrics.IncreaseGenUrlRequests(inputData.User)
	shortUrl := this.redis.Get(inputData.Url)
	if shortUrl != "" {
		return shortUrl, nil
	}
	urlRecord, err := this.urlRepository.FindLongUrl(inputData.Url)
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	// Found url and not exprire
	if err != gorm.ErrRecordNotFound && urlRecord.ExpireTime.After(time.Now()) {
		return urlRecord.ShortUrl, nil
	}
	// insert DB
	channelRepo := make(chan struct{})
	channelRedis := make(chan struct{})
	newID := this.GetNextID()
	newShortUrl := encryption.GenerateShortLink(newID)
	var errRepo, errRedis error
	go func() {
		errRepo = this.urlRepository.InsertUrl(newID, inputData.User, newShortUrl, inputData.Url, time.Now().AddDate(0, 0, 3))
		channelRepo <- struct{}{}
	}()
	go func() {
		errRedis = this.redis.SetJSON(newShortUrl, inputData, 24*time.Hour)
		if errRedis != nil {
			channelRedis <- struct{}{}
			return
		}
		errRedis = this.redis.Set(inputData.Url, newShortUrl, 24*time.Hour)
		channelRedis <- struct{}{}
	}()
	go this.metrics.ResetGetUrlRequests(newShortUrl, inputData.User)
	<-channelRepo
	<-channelRedis
	if errRepo != nil {
		return "", errRepo
	}
	if errRedis != nil {
		return "", errRedis
	}
	return newShortUrl, nil
}
