package databases

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"talkspace-api/app/configs"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	redisClient *redis.Client
	once        sync.Once
	ctx         = context.Background()
)

func ConnectRedis() *redis.Client {
	once.Do(func() {
		config, err := configs.LoadConfig()
		if err != nil {
			logrus.Fatalf("failed to load Redis configuration: %v", err)
			return
		}

		redisDB, err := strconv.Atoi(config.REDIS.REDIS_DB)
		if err != nil {
			redisDB = 1
		}

		client := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s",
				config.REDIS.REDIS_HOST,
				config.REDIS.REDIS_PORT,
			),
			Password: config.REDIS.REDIS_PASS,
			DB:       redisDB,
		})

		_, err = client.Ping(ctx).Result()
		if err != nil {
			logrus.Errorf("failed to connect to Redis: %v", err)
			return
		}

		logrus.Info("connected to Redis")
		redisClient = client
	})

	return redisClient
}

func SetToken(key string, value string, expiration time.Duration) error {
	err := redisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		logrus.Errorf("failed to set token in Redis: %v", err)
	}
	return err
}

func GetToken(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		logrus.Errorf("failed to get token from Redis: %v", err)
	}
	return value, err
}
