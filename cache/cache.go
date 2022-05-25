package cache

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	connection *redis.Client
}

func (this *Cache) Connect() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		fmt.Println(err.Error())

	} else {
		fmt.Println("Success connect to redis")
		this.connection = client
	}
}
func (this *Cache) Get(key string) (string, error) {
	return this.connection.Get(this.connection.Context(), key).Result()
}
func (this *Cache) Set(key string, value string, TTL int) error {
	log.Println("abc", key, value, TTL)
	fmt.Println("abc", key, value, TTL)
	return this.connection.Set(this.connection.Context(), key, value, time.Duration(TTL)*time.Hour).Err()
}
func (this *Cache) Flush() error {
	return this.connection.FlushDB(this.connection.Context()).Err()
}
func (this *Cache) Increase(key string) int64 {
	return this.connection.Incr(this.connection.Context(), key).Val()
}
