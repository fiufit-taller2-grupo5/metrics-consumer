package redis_client

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

const redisPort = 6379

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	host := getHostFromEnvironment()
	address := fmt.Sprintf("%s:%d", host, redisPort)
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr: address,
		}),
	}
}

func (rc *RedisClient) GetAllFromList(key string) ([]string, error) {
	result, err := rc.client.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rc *RedisClient) RemoveAllFromList(key string) error {
	_, err := rc.client.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return err
}

func getHostFromEnvironment() string {
	if os.Getenv("environment") == "production" {
		return "redis-service"
	} else {
		return "localhost"
	}
}
