package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var Client *redis.Client

type Config struct {
	Host           string
	Port           int
	User           string
	Password       string
	Db             int
	ConnectTimeout int
	ReadTimeout    int
	WriteTimeout   int
	PoolSize       int
}

func GetConn(config Config) (*redis.Client, error) {
	if Client != nil && Client.Ping(GetCtx()).Err() == nil {
		return Client, nil
	}

	Client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username:     config.User,
		Password:     config.Password,
		DB:           config.Db,
		DialTimeout:  time.Duration(config.ConnectTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
		PoolSize:     config.PoolSize,
	})
	if err := Client.Ping(GetCtx()).Err(); err != nil {
		return nil, err
	}
	return Client, nil
}

func GetCtx() context.Context {
	return context.Background()
}