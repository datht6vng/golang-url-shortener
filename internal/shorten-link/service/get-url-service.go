package service

import (
	"time"
	"trueid-shorten-link/internal/shorten-link/repository"
	"trueid-shorten-link/package/metrics"
	_redis "trueid-shorten-link/package/redis"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GetURLService struct {
	urlRepository *repository.URLRepository
	redis         *_redis.Redis
	metrics       *metrics.Metrics
}

func (this *GetURLService) Init(urlRepository *repository.URLRepository, redis *_redis.Redis, metrics *metrics.Metrics) *GetURLService {
	this.urlRepository = urlRepository
	this.redis = redis
	this.metrics = metrics
	return this
}

func (this *GetURLService) GetURL(url string) (string, error) {
	cacheURL := new(URLData)
	err := this.redis.GetJSON("url:"+url, cacheURL)
	if err != nil {
		return "", err
	}
	if cacheURL.URL != "" {
		go this.metrics.IncreaseGetURLRequests(url, cacheURL.ClientID)
		return cacheURL.URL, nil
	}
	urlRecord, err := this.urlRepository.FindByShortURL(url)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", &fiber.Error{
				Code:    404,
				Message: "URL not found!",
			}
		}
		return "", err
	}
	if urlRecord.ExpireTime.Before(time.Now()) {
		return "", &fiber.Error{
			Code:    410,
			Message: "URL has been expired!",
		}
	}
	err = this.redis.SetJSON("url:"+urlRecord.ShortURL, URLData{URL: urlRecord.ShortURL, ClientID: urlRecord.ClientID}, 24*time.Hour)
	err = this.redis.Set("url:"+urlRecord.LongURL, urlRecord.ShortURL, 24*time.Hour)
	if err != nil {
		return "", err
	}
	go this.metrics.IncreaseGetURLRequests(url, urlRecord.ClientID)
	return urlRecord.LongURL, nil
}
