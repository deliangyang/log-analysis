package main

import (
	"github.com/hpcloud/tail"
	"fmt"
	"regexp"
	"log-analysis/log"
	"github.com/json-iterator/go"
)

var pattern, _ = regexp.Compile("{[^\n]+")
var pathPattern, _ = regexp.Compile(`(^[^\?]+)`)

func main() {

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
