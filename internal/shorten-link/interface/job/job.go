package job

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
	"trueid-shorten-link/internal/shorten-link/repository"
	"trueid-shorten-link/internal/shorten-link/service"
	_redis "trueid-shorten-link/package/redis"

	"github.com/robfig/cron/v3"
)

type Job struct {
	urlRepository             *repository.URLRepository
	generateCounterRepository *repository.GenerateCounterRepository
	redis                     *_redis.Redis
	cron                      *cron.Cron
}

func (this *Job) Init(urlRepository *repository.URLRepository, generateCounterRepository *repository.GenerateCounterRepository, redis *_redis.Redis) *Job {
	this.urlRepository = urlRepository
	this.generateCounterRepository = generateCounterRepository
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
func (this *Job) DeleteExpireURL() {
	this.urlRepository.DeleteExpiredURL()
	log.Printf("Delete expire URL!")
}
func (this *Job) ResetMaxID() {
	maxID, _ := this.urlRepository.GetMaxID()
	this.redis.Set("CurrentID", maxID, 0)
	log.Printf("Reset max ID!")
}

// for backup version of counter, current solution use Redis and query directly from url table
func (this *Job) BackupGenCounter() {
	keys := this.redis.Keys("gen-counter:*")
	for _, key := range keys {
		generateCounterStr := this.redis.Get(key)
		this.redis.Set(key, "0", 24*time.Hour)
		if generateCounterStr == "" {
			continue
		}
		generateCounter, _ := strconv.ParseInt(generateCounterStr, 10, 64)
		key = strings.Split(key, "gen-counter:")[1]
		counterKey := new(service.GenerateCounterKey)
		_ = json.Unmarshal([]byte(key), counterKey)
		err := this.generateCounterRepository.Insert(counterKey.ClientID, counterKey.CreateDate, generateCounter)
		if err != nil {
			this.generateCounterRepository.Update(counterKey.ClientID, counterKey.CreateDate, generateCounter)
		}

	}
}
