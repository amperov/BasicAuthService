package redis

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type ClientRedis struct {
	red *redis.Client
}

func GetRedisClient() *ClientRedis {
	RedisOptions := redis.Options{
		Addr:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"),
		DB:       0,
	}
	Client := redis.NewClient(&RedisOptions)

	ping := Client.Ping()
	if ping.Err() != nil {
		logrus.Fatal(ping.Err())
	}
	return &ClientRedis{red: Client}
}

func (c *ClientRedis) InsertAccessToken(ctx context.Context, AccessCode, AccessToken string) error {
	set := c.red.Set(AccessCode, AccessToken, time.Hour*3)
	if set.Err() != nil {
		return set.Err()
	}
	return nil
}

func (c *ClientRedis) GetAccessToken(ctx context.Context, AccessCode string) (string, error) {
	get := c.red.Get(AccessCode)
	result, err := get.Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
