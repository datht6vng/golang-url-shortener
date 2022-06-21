package service

import (
	"time"
	"trueid-shorten-link/internal/shorten-link/entity"
	"trueid-shorten-link/internal/shorten-link/repository"
	"trueid-shorten-link/package/encryption"
	logger "trueid-shorten-link/package/log"
	"trueid-shorten-link/package/metrics"
	_redis "trueid-shorten-link/package/redis"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

type GenerateURLService interface {
	GenerateURL(url string, clientID string) (string, error)
}

type generateURLService struct {
	urlRepository *repository.URLRepository
	redis         *_redis.Redis
	metrics       *metrics.Metrics
	timeFormat    string
}

func NewGenerateURLService(urlRepository *repository.URLRepository, redis *_redis.Redis, metrics *metrics.Metrics) GenerateURLService {
	return &generateURLService{
		urlRepository: urlRepository,
		redis:         redis,
		metrics:       metrics,
		timeFormat:    "2006-01-02",
	}
}

func (s *generateURLService) GetNextID() int64 {
	if s.redis.Get("CurrentID") == "" {
		pipe := s.redis.TxPipeline()
		var resultID int64
		s.redis.Watch(func(tx *redis.Tx) error {
			maxID, err := s.urlRepository.GetMaxID()
			if err != nil {
				logger.GetLog().Errorf("Error reset ID %v", err.Error())
				return err
			}
			pipe.Set("CurrentID", maxID, -1)
			cmd := pipe.Incr("CurrentID")
			if _, err := pipe.Exec(); err != nil && err != redis.Nil {
				logger.GetLog().Errorf("Error reset ID %v", err.Error())
				return err
			}
			resultID = cmd.Val()
			return nil
		}, "CurrentID")
		return resultID
	}
	return s.redis.Incr("CurrentID")
}

func (s *generateURLService) GenerateURL(url string, clientID string) (string, error) {
	logger := logger.GetLog()
	go s.metrics.IncreaseGenURLRequests(clientID)
	shortURL := s.redis.HGet("LONG_URLS", url)
	if shortURL != "" {
		go func() {
			err := s.redis.HSetTTL("SHORT_URLS", shortURL, entity.URLData{URL: url, ClientID: clientID}, 24*time.Hour)
			if err != nil {
				logger.Errorf("Set cache fail: %v", err)
				return
			}
		}()
		go func() {
			err := s.redis.HSetTTL("LONG_URLS", url, shortURL, 24*time.Hour)
			if err != nil {
				logger.Errorf("Set cache fail: %v", err)
			}
		}()
		return shortURL, nil
	}
	urlRecord, err := s.urlRepository.FindByLongURL(url)
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	// Found url and not exprire
	if err != gorm.ErrRecordNotFound && urlRecord.ExpireTime.After(time.Now()) {
		go func() {
			err := s.redis.HSetTTL("SHORT_URLS", urlRecord.ShortURL, entity.URLData{URL: url, ClientID: clientID}, 24*time.Hour)
			if err != nil {
				logger.Errorf("Set cache fail: %v", err)
			}
		}()
		go func() {
			err := s.redis.HSetTTL("LONG_URLS", url, urlRecord.ShortURL, 24*time.Hour)
			if err != nil {
				logger.Errorf("Set cache fail: %v", err)
			}
		}()
		return urlRecord.ShortURL, nil
	}
	// insert DB
	newID := s.GetNextID()

	newShortURL := encryption.GenerateShortLink(newID)
	go func() {
		err := s.urlRepository.InsertURL(newID, clientID, newShortURL, url, time.Now().AddDate(0, 0, 3))

		if err != nil {
			logger.Errorf("Insert DB fail: %v", err)
		}
	}()
	go func() {
		err := s.redis.HSetTTL("SHORT_URLS", newShortURL, entity.URLData{URL: url, ClientID: clientID}, 24*time.Hour)
		if err != nil {
			logger.Errorf("Set cache fail: %v", err)
		}
	}()
	go func() {
		err := s.redis.HSetTTL("LONG_URLS", url, newShortURL, 24*time.Hour)
		if err != nil {
			logger.Errorf("Set cache fail: %v", err)
		}
	}()
	go s.metrics.ResetGetURLRequests(newShortURL, clientID)
	return newShortURL, nil
}
