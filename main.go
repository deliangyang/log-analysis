package main

import (
	"github.com/hpcloud/tail"
	"fmt"
	"regexp"
	"log-analysis/log"
	"github.com/json-iterator/go"
	"github.com/go-redis/redis"
	"log-analysis/request"
	"os"
)

var pattern, _ = regexp.Compile("{[^\n]+")
var pathPattern, _ = regexp.Compile(`(^[^\?]+)`)

var redisOptions = &redis.Options{
	Addr: "127.0.0.1",
}

func main() {

	local := request.Location{}
	local.Query("58.17.200.100")
	os.Exit(1)
	//redisCli := redis.NewClient(redisOptions)
/*	db, err := gorm.Open("mysql", "root:www.ydl.com@/test?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	db.AutoMigrate(&model.Product{})

	db.NewRecord(&model.Product{
		LanguageCode: "1ddd",
		Code: "xxxxx",
		Name: "test",
	})*/

	t, err := tail.TailFile("data/api-2018-09-04.log", tail.Config{Follow: true})
	if err != nil {
		panic(err)
	}
	for line := range t.Lines {
		date := line.Text[1:20]
		fmt.Println(date)
		content := pattern.FindString(line.Text)
		//fmt.Println(content)
		var httpRequest log.HttpRequestLog
		jsoniter.Unmarshal([]byte(content), &httpRequest)
		httpRequest.Path = pathPattern.FindString(httpRequest.Path)
		if httpRequest.Path == "/api/search" {
			fmt.Println(string(httpRequest.Request["keywords"]))
		}
	}
}
