package main

import (
	"github.com/go-redis/redis"
	"log-analysis/log"
	"sync"
)

var redisOptions = &redis.Options{
	Addr: "127.0.0.1:6379",
	DB: 2,
}

func main() {
	redisCli := redis.NewClient(redisOptions)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Path = "data/"
		fileLog := log.FileLog{
			Filename: "data/",
			Line: 1,
			RedisClient: redisCli,
			Date: "2018-09-04",
		}
		fileLog.TailFile("api")
	}()
	wg.Wait()
}
