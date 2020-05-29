package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis-server",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()

	fmt.Println(pong, err)
}
