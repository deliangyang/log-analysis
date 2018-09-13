package main

import (
	"github.com/go-redis/redis"
	"log-analysis/log"
	"sync"
)

var redisOptions = &redis.Options{
	Addr: "www.ydl.com:6379",
}

func main() {
	redisCli := redis.NewClient(redisOptions)
/*	db, err := gorm.Open("mysql", "root:www.ydl.com@/test?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	db.AutoMigrate(&model.Product{})

	db.NewRecord(&model.Product{
		LanguageCode: "1ddd",
		Code: "xxxxx",
		Name: "test",
	})*/
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
	//time.Sleep(time.Millisecond * 10000)
	/*go func() {
		fileLog := log.FileLog{
			Filename: "data/",
			Line: 1,
			RedisClient: redisCli,
			Date: "2018-09-04",
		}
		fileLog.TailFile("api")
	}()*/
}
