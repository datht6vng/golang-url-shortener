package service

import (
	"trueid-shorten-link/internal/shorten-link/repository"
	"trueid-shorten-link/package/metrics"
	_redis "trueid-shorten-link/package/redis"
)

type generateURLServiceVer2 struct {
	urlRepository *repository.URLRepository
	redis         *_redis.Redis
	metrics       *metrics.Metrics
	timeFormat    string
}

func NewGenerateURLServiceVer2(urlRepository *repository.URLRepository, redis *_redis.Redis, metrics *metrics.Metrics) GenerateURLService {
	return &generateURLServiceVer2{
		urlRepository: urlRepository,
		redis:         redis,
		metrics:       metrics,
		timeFormat:    "2006-01-02",
	}
}

func (s *generateURLServiceVer2) GenerateURL(url string, clientID string) (string, error) {
	go s.metrics.IncreaseGenURLRequests(clientID)

	shortURL := s.redis.Get("URL:" + url)
	if shortURL != "" {
		return shortURL, nil
	}

	panic("Version 2")
}
