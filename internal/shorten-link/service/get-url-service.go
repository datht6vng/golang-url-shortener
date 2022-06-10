package service

import (
	"time"
	"trueid-shorten-link/internal/shorten-link/repository"
	"trueid-shorten-link/package/metrics"
	_redis "trueid-shorten-link/package/redis"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GetUrlService struct {
	urlRepository *repository.UrlRepository
	redis         *_redis.Redis
	metrics       *metrics.Metrics
}

func (this *GetUrlService) Init(urlRepository *repository.UrlRepository, redis *_redis.Redis, metrics *metrics.Metrics) *GetUrlService {
	this.urlRepository = urlRepository
	this.redis = redis
	this.metrics = metrics
	return this
}

func (this *GetUrlService) GetUrl(url string) (string, error) {
	cacheUrl := new(UrlData)
	err := this.redis.GetJSON(url, cacheUrl)
	if err != nil {
		return "", err
	}
	if cacheUrl.Url != "" {
		go this.metrics.IncreaseGetUrlRequests(url, cacheUrl.User)
		return cacheUrl.Url, nil
	}
	urlRecord, err := this.urlRepository.FindShortUrl(url)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", &fiber.Error{
				Code:    404,
				Message: "Url not found!",
			}
		}
		return "", err
	}
	if urlRecord.ExpireTime.Before(time.Now()) {
		return "", &fiber.Error{
			Code:    410,
			Message: "Url has been expired!",
		}
	}
	err = this.redis.SetJSON(urlRecord.ShortUrl, UrlData{Url: urlRecord.ShortUrl, User: urlRecord.User}, 24*time.Hour)
	err = this.redis.Set(urlRecord.LongUrl, urlRecord.ShortUrl, 24*time.Hour)
	if err != nil {
		return "", err
	}
	go this.metrics.IncreaseGetUrlRequests(url, urlRecord.User)
	return urlRecord.LongUrl, nil
}
