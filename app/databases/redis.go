package databases

import (
	"fmt"
	"strconv"
	"talkspace-api/app/configs"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func ConnectRedis() *redis.Client {
	config, err := configs.LoadConfig()
	if err != nil {
		logrus.Fatalf("failed to load Redis configuration: %v", err)
		return nil
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

	_, err = client.Ping().Result()
	if err != nil {
		logrus.Errorf("failed to connect to Redis: %v", err)
		return nil
	}

	logrus.Info("connected to Redis")
	
	return client
}