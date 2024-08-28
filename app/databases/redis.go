package databases

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
	"talkspace-api/app/configs"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	rdb  *redis.Client
	once sync.Once
	ctx  = context.Background()
)

func ConnectRedis() *redis.Client {
	once.Do(func() {
		config, err := configs.LoadConfig()
		if err != nil {
			logrus.Fatalf("failed to load Redis configuration: %v", err)
		}

		redisDB, err := strconv.Atoi(config.REDIS.REDIS_DB)
		if err != nil {
			redisDB = 1
		}

		client := redis.NewClient(&redis.Options{
			Addr:        fmt.Sprintf("%s:%s", config.REDIS.REDIS_HOST, config.REDIS.REDIS_PORT),
			Password:    config.REDIS.REDIS_PASS,
			DB:          redisDB,
			MaxRetries:  3,
			DialTimeout: 10 * time.Second,
			ReadTimeout: 10 * time.Second,
		})

		_, err = client.Ping(ctx).Result()
		if err != nil {
			logrus.Fatalf("failed to connect to Redis: %v", err)
		}

		logrus.Info("connected to Redis")
		rdb = client
	})

	return rdb
}
