package service

import (
	"trueid-shorten-link/internal/shorten-link/repository"
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
func (this *ValidateAPIKeyService) ValidateAPIKey(apiKey string) (string, error) {
	cacheClientID := this.redis.Get("API-KEY:" + apiKey)
	if cacheClientID != "" {
		return cacheClientID, nil
	}
	clientIDRecord, err := this.clientRepository.FindByAPIKey(apiKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", &fiber.Error{
				Code:    401,
				Message: "API Key not found!",
			}
		}
		return "", err
	}
	return clientIDRecord.ClientID, nil
}
