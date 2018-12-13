package main

import (
	"flag"
	"os"
	"fmt"
	"bufio"
	"io"
	"time"
	"sync"
	"net/http"
	"io/ioutil"
	"regexp"
)

var (
	pattern, _ = regexp.Compile(`反向链接：</span>(\d+)</li>`)
)

func main() {
	var filename string

	flag.StringVar(&filename, "filename", "", "enter your filename")
	flag.Parse()

	if filename == "" {
		os.Exit(-1)
	}

	waitGroup := sync.WaitGroup{}
	domainList := make(chan string)
	go readFile(filename, domainList)


	for i := 0; i < 1; i++ {
		waitGroup.Add(1)
		go func(domainList chan string) {
			for {
				select {
				case domain := <-domainList:
					lines := fetchPage(domain)
					println(lines)
					fmt.Println(lines)
				case <-timeoutAfter(time.Second * 4):
					waitGroup.Done()
				}
			}
		}(domainList)
	}

	waitGroup.Wait()
}

func timeoutAfter(d time.Duration) chan int {
	q := make(chan int, 1)

	time.AfterFunc(d, func() {
		q <- 1
	})

	return q
}

func fetchPage(url string) string {
	targetUrl := "http://alexa.chinaz.com/?domain=" + url
	req, err := http.Get(targetUrl)
	if err != nil {
		return ""
	}
	body, _ := ioutil.ReadAll(req.Body)
	value := pattern.FindSubmatch(body)
	linkNum := "未获取到"
	if len(value) >= 2 {
		linkNum = string(value[1])
	}
	time.Sleep(time.Second * 3)
	return url + "\t" + linkNum
}

func readFile(filename string, urlChannel chan string) {
	fi, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		urlChannel <- string(a)
	}
}
