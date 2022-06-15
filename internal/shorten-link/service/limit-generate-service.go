package service

import (
	"fmt"
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
	counterKey := clientID + time.Now().Format(this.timeFormat)
	if this.redis.Get(counterKey) == "" {
		// Create pipline to reset
		fmt.Println("Run in thiss")
		pipe := this.redis.TxPipeline()
		this.redis.Watch(func(tx *redis.Tx) error {
			limit, _ := this.urlRepository.CountLinkGenerated(clientID)
			pipe.Set(counterKey, limit, 24*time.Hour)
			if _, err := pipe.Exec(); err != nil && err != redis.Nil {
				return err
			}
			return nil
		}, counterKey)
	}
	// Check client limit gen
	counter, err := strconv.ParseInt(this.redis.Get(counterKey), 10, 64)
	if err != nil {
		return err
	}
	if counter > limit {
		return &fiber.Error{
			Code:    429,
			Message: "Reach limit of link generation!",
		}
	}
	return nil
}