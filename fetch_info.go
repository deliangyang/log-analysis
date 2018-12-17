package main

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"fmt"
	"time"
)

type WebSiteInfo struct {
	Domain string
	Ip string
	Title string
	KeyWords string
	Description string
}

var (
	Client = &http.Client{Timeout: time.Second * 3}
	TitlePattern, _ = regexp.Compile(`<title>([^<]+)</title>`)
	KeywordPattern, _ = regexp.Compile(`<meta[^>]+name="keywords"[^>]+content="([^"]+)"`)
	DescriptionPattern, _ = regexp.Compile(`<meta[^>]+name="description"[^>]+content="([^"]+)"`)
)

func main() {
	urls := []string {"http://www.baidu.com", "http://www.admin5.com", "http://www.alipay.com", "https://blog.csdn.net"}
	for _, url := range urls {
		websiteInfo, _ := getHtml(url)
		fmt.Println(websiteInfo.Title)
		fmt.Println(websiteInfo.KeyWords)
		fmt.Println(websiteInfo.Description)
	}
}

func getHtml(domain string) (WebSiteInfo, error) {
	req, err := http.NewRequest("GET", domain, nil)
	if err != nil {
		fmt.Println(err)
	}
	rsp, err := Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return WebSiteInfo{}, nil
	}
	defer rsp.Body.Close()

	content, _ := ioutil.ReadAll(rsp.Body)
	return parseInfo(content), nil
}

func parseInfo(content []byte) WebSiteInfo {
	title := TitlePattern.FindSubmatch(content)
	websiteInfo := WebSiteInfo{}
	if len(title) >= 2 {
		websiteInfo.Title = string(title[1])
	}

	keywords := KeywordPattern.FindSubmatch(content)
	if len(keywords) >= 2 {
		websiteInfo.KeyWords = string(keywords[1])
	}

	description := DescriptionPattern.FindSubmatch(content)
	if len(description) >= 2 {
		websiteInfo.Description = string(description[1])
	}
	return websiteInfo
}