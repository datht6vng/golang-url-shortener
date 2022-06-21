package service

import (
	"fmt"
	"trueid-shorten-link/internal/shorten-link/entity"
	"trueid-shorten-link/internal/shorten-link/repository"
	"trueid-shorten-link/package/redis"
)

type ClientService struct {
	repo  *repository.ClientRepository
	cache *redis.Redis
}

func NewClientService(repo *repository.ClientRepository, cache *redis.Redis) *ClientService {
	return &ClientService{
		repo:  repo,
		cache: cache,
	}
}

func (s *ClientService) UpdateClient(updateRequest *entity.UpdateClientRequest) error {
	client, err := s.repo.FindByID(updateRequest.ClientId)
	if err != nil {
		return err
	}
	switch updateRequest.UpdateType {
	case string(entity.UpdateLimit):
		client.MaxLink = updateRequest.Limit
	default:
		return fmt.Errorf("unsupported update Type")
	}

	if err := s.repo.Update(client); err != nil {
		return err
	}
	s.cache.HDel("API-KEY:"+client.APIKey, "client")
	return nil
}
