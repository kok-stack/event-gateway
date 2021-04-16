package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/kok-stack/event-gateway/pkg/config"
)

func Conn(config *config.ApplicationConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "192.168.212.213:6379",
		//Addr:     config.Redis.Addr,
		//Password: config.Redis.Password, // no password set
		DB: 0, // use default DB
	})
	return rdb
}
