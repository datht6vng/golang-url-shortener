package job

import (
	"log"
	"trueid-shorten-link/internal/shorten-link/repository"
	_redis "trueid-shorten-link/package/redis"

	"github.com/robfig/cron/v3"
)

type Job struct {
	urlRepository *repository.UrlRepository
	redis         *_redis.Redis
	cron          *cron.Cron
}

func (this *Job) Init(urlRepository *repository.UrlRepository, redis *_redis.Redis) *Job {
	this.urlRepository = urlRepository
	this.redis = redis
	this.cron = cron.New()
	return this
}
func (this *Job) CreateCronJob(interval string, cronJobs ...func()) {
	for _, cronJob := range cronJobs {
		this.cron.AddFunc(interval, cronJob)
	}
	this.cron.Start()
}
func (this *Job) DeleteExpireUrl() {
	this.urlRepository.DeleteExpiredUrl()
	log.Printf("Delete expire URL!")
}
func (this *Job) ResetMaxID() {
	maxID, _ := this.urlRepository.GetMaxID()
	this.redis.Set("CurrentID", maxID, 0)
	log.Printf("Reset max ID!")
}
