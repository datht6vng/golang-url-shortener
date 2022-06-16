package service

import (
	"strconv"
	"time"
	"trueid-shorten-link/internal/shorten-link/repository"
	_redis "trueid-shorten-link/package/redis"

	"github.com/go-redis/redis/v7"
	"github.com/gofiber/fiber/v2"
)

type LimitGenerateService struct {
	urlRepository *repository.URLRepository
	redis         *_redis.Redis
	timeFormat    string
}

func (this *LimitGenerateService) Init(urlRepository *repository.URLRepository, redis *_redis.Redis) *LimitGenerateService {
	this.urlRepository = urlRepository
	this.redis = redis
	this.timeFormat = "2006-01-02"
	return this
}
func (this *LimitGenerateService) LimitGenerate(clientID string, limit int64) error {
	counterKey := "LinkCounter:" + time.Now().Format(this.timeFormat) + "|" + clientID
	result := this.redis.Get(counterKey)
	var counter int64
	if result == "" {
		// Create pipline to reset
		pipe := this.redis.TxPipeline()
		this.redis.Watch(func(tx *redis.Tx) error {
			counter, _ = this.urlRepository.CountLinkGenerated(clientID)
			pipe.Set(counterKey, counter, 24*time.Hour)
			if _, err := pipe.Exec(); err != nil && err != redis.Nil {
				return err
			}
			return nil
		}, counterKey)
	} else {
		var err error
		counter, err = strconv.ParseInt(result, 10, 64)
		if err != nil {
			return err
		}
	}
	if counter > limit {
		return &fiber.Error{
			Code:    429,
			Message: "Reach limit of link generation!",
		}
	}
	go this.redis.Incr(counterKey)
	return nil
}
