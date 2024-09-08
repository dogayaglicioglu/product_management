package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

var RedisC Redis

func ConnectRedis(c chan bool) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in loading env. varibale %v", err)
	}
	redisPwd := os.Getenv("REDIS_PWD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: redisPwd,
		DB:       0, // to use default DB
	})

	if err := checkRedisConn(rdb); err != nil {
		log.Fatal("Redis conn. failed %v", err)
	}
	fmt.Println("Redis connection successful!")
	RedisC = Redis{
		Client: rdb,
	}
	c <- true

}

func checkRedisConn(rdb *redis.Client) error {
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to ping Redis server: %w", err)
	}
	return nil
}

func GetRedClient() *redis.Client {
	return RedisC.Client
}
