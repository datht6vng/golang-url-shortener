package service

import (
	"encoding/json"
	"time"
	"trueid-shorten-link/internal/shorten-link/repository"
	"trueid-shorten-link/package/encryption"
	"trueid-shorten-link/package/metrics"
	_redis "trueid-shorten-link/package/redis"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

type GenerateURLService struct {
	urlRepository *repository.URLRepository
	redis         *_redis.Redis
	metrics       *metrics.Metrics
	timeFormat    string
}

func (this *GenerateURLService) Init(urlRepository *repository.URLRepository, redis *_redis.Redis, metrics *metrics.Metrics) *GenerateURLService {
	this.urlRepository = urlRepository
	this.redis = redis
	this.metrics = metrics
	this.timeFormat = "2006-01-02"
	return this
}

func (this *GenerateURLService) GetNextID() int64 {
	if currentID := this.redis.Get("CurrentID"); currentID == "" {
		pipe := this.redis.TxPipeline()
		var resultID int64
		this.redis.Watch(func(tx *redis.Tx) error {
			maxID, _ := this.urlRepository.GetMaxID()
			pipe.Set("CurrentID", maxID, -1)
			cmd := pipe.Incr("CurrentID")
			if _, err := pipe.Exec(); err != nil && err != redis.Nil {
				return err
			}
			resultID = cmd.Val()
			return nil
		}, "CurrentID")
		return resultID
	}
	return this.redis.Incr("CurrentID")
}

func (this *GenerateURLService) GenerateURL(url string, clientID string) (string, error) {
	go this.metrics.IncreaseGenURLRequests(clientID)
	counterKey, _ := json.Marshal(GenerateCounterKey{
		ClientID:   clientID,
		CreateDate: time.Now().Format(this.timeFormat),
	})
	go func() {
		this.redis.Incr("gen-counter:" + string(counterKey))
		this.redis.Expire("gen-counter:"+string(counterKey), 24*time.Hour)
	}()
	shortURL := this.redis.Get("url:" + url)
	if shortURL != "" {
		return shortURL, nil
	}
	urlRecord, err := this.urlRepository.FindByLongURL(url)
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	// Found url and not exprire
	if err != gorm.ErrRecordNotFound && urlRecord.ExpireTime.After(time.Now()) {
		return urlRecord.ShortURL, nil
	}
	// insert DB
	channelRepo := make(chan struct{})
	channelRedis := make(chan struct{})
	newID := this.GetNextID()
	newShortURL := encryption.GenerateShortLink(newID)
	var errRepo, errRedis error
	go func() {
		errRepo = this.urlRepository.InsertURL(newID, clientID, newShortURL, url, time.Now().AddDate(0, 0, 3))
		channelRepo <- struct{}{}
	}()
	go func() {
		errRedis = this.redis.SetJSON("url:"+newShortURL, URLData{URL: url, ClientID: clientID}, 24*time.Hour)
		if errRedis != nil {
			channelRedis <- struct{}{}
			return
		}
		errRedis = this.redis.Set("url:"+url, newShortURL, 24*time.Hour)
		channelRedis <- struct{}{}
	}()
	go this.metrics.ResetGetURLRequests(newShortURL, clientID)
	<-channelRepo
	<-channelRedis
	if errRepo != nil {
		return "", errRepo
	}
	if errRedis != nil {
		return "", errRedis
	}
	return newShortURL, nil
}
