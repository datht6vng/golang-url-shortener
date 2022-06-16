package service

import (
	"time"
	"trueid-shorten-link/internal/shorten-link/repository"
	"trueid-shorten-link/package/model"
	_redis "trueid-shorten-link/package/redis"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ValidateAPIKeyService struct {
	clientRepository *repository.ClientRepository
	redis            *_redis.Redis
}

func (this *ValidateAPIKeyService) Init(clientRepository *repository.ClientRepository, redis *_redis.Redis) *ValidateAPIKeyService {
	this.clientRepository = clientRepository
	this.redis = redis
	return this
}
func (this *ValidateAPIKeyService) ValidateAPIKey(apiKey string) (string, int64, error) {
	var cacheClient = new(model.Client)
	this.redis.GetJSON("API-KEY:"+apiKey, cacheClient)
	if cacheClient.ClientID != "" {
		return cacheClient.ClientID, cacheClient.MaxLink, nil
	}
	clientRecord, err := this.clientRepository.FindByAPIKey(apiKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", 0, &fiber.Error{
				Code:    401,
				Message: "API Key not found!",
			}
		}
		return "", 0, err
	}
	this.redis.SetJSON("API-KEY:"+apiKey, clientRecord, 24*time.Hour)
	return clientRecord.ClientID, clientRecord.MaxLink, nil
}
