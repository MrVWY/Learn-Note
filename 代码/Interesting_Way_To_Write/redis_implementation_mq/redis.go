package redis_implementation_mq

import "github.com/go-redis/redis/v8"

type MQ struct {
	Name string
	*redis.Client
}

func NewMQ(addr, name string) *MQ {

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return &MQ{
		Name:   name,
		Client: rdb,
	}
}
