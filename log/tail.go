package log

import (
	"github.com/json-iterator/go"
	"fmt"
	"github.com/hpcloud/tail"
	"regexp"
	"github.com/go-redis/redis"
	"io/ioutil"
)

var Path = "/var/wwww"
var pattern, _ = regexp.Compile("{[^\n]+")
var pathPattern, _ = regexp.Compile(`(^[a-zA-Z/]+/?)`)

type FileLog struct {
	Filename string
	Date string
	Line int
	RedisClient redis.Cmdable
}

func (fileLog *FileLog) TailFile(channel string)  {
	fmt.Println(Path + channel + "-" + fileLog.Date + ".log")
	t, err := tail.TailFile(Path + channel + "-" + fileLog.Date + ".log", tail.Config{Follow: true})
	if err != nil {
		panic(err)
	}
	for line := range t.Lines {
		ioutil.WriteFile("test", []byte(line.Text), 0644)
		content := pattern.FindString(line.Text)
		var httpRequest HttpRequestLog
		contentByte := []byte(content)
		jsoniter.Unmarshal(contentByte, &httpRequest)
		userId := jsoniter.Get(contentByte, "userId").ToInt()
		httpRequest.Path = pathPattern.FindString(httpRequest.Path)
		fileLog.RedisClient.HIncrBy(fileLog.Date, httpRequest.Path, 1).Result()
		if httpRequest.Path == "/api/search" {
			fmt.Println(string(httpRequest.Request["keywords"]))
			fileLog.RedisClient.HIncrBy("keyword:" + fileLog.Date, string(httpRequest.Request["keywords"]), 1)
		} else if httpRequest.Path == "/api/products/" {
			response := jsoniter.Get(contentByte, "response").ToString()
			responseByte := []byte(response)
			id := jsoniter.Get(responseByte, "id").ToInt()
			title := jsoniter.Get(responseByte, "title").ToString()
			amount := jsoniter.Get(responseByte, "amount").ToInt()
			fmt.Println(userId, id, title, amount)
		}
	}
}
