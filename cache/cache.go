package cache

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	connection *redis.Client
}

func (this *Cache) Connect() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "127.0.0.1"
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
