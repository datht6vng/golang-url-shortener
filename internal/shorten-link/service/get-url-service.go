package service

import (
	"time"
	"trueid-shorten-link/internal/shorten-link/entity"
	"trueid-shorten-link/internal/shorten-link/repository"
	logger "trueid-shorten-link/package/log"
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
	logger := logger.GetLog()
	cacheURL := new(entity.URLData)
	err := this.redis.HGetJSON("SHORT_URLS", url, cacheURL)
	if err != nil {
		return "", err
	}

	if cacheURL.URL != "" {
		go func() {
			err := this.redis.HSetTTL("SHORT_URLS", url, entity.URLData{URL: cacheURL.URL, ClientID: cacheURL.ClientID}, 24*time.Hour)
			if err != nil {
				logger.Errorf("Set cache fail %v", err.Error())
			}
		}()
		go func() {
			err = this.redis.HSetTTL("LONG_URLS", cacheURL.URL, url, 24*time.Hour)
			if err != nil {
				logger.Errorf("Set cache fail %v", err.Error())
			}
		}()
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
	go func() {
		err := this.redis.HSetTTL("SHORT_URLS", urlRecord.ShortURL, entity.URLData{URL: urlRecord.LongURL, ClientID: urlRecord.ClientID}, 24*time.Hour)
		if err != nil {
			logger.Errorf("Set cache fail %v", err.Error())
		}
	}()
	go func() {
		err = this.redis.HSetTTL("LONG_URLS", urlRecord.LongURL, urlRecord.ShortURL, 24*time.Hour)
		if err != nil {
			logger.Errorf("Set cache fail %v", err.Error())
		}
	}()
	go this.metrics.IncreaseGetURLRequests(url, urlRecord.ClientID)
	return urlRecord.LongURL, nil
}
