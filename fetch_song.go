package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"io"
	"bufio"
	"sync"
	"strings"
	"time"
	URL "net/url"
	"sync/atomic"
)

var (
	waitGroup = sync.WaitGroup{}

	count uint64 = 0
)

func main() {
	var thread int
	var isDebug bool
	var filename string

	flag.IntVar(&thread, "thread", 5, "work thread number")
	flag.BoolVar(&isDebug, "debug", true, "default mode is debug")
	flag.StringVar(&filename, "filename", "t3.log", "fetch url filename")
	flag.Parse()

	if filename == "" {
		panic("error filename is nil")
	}

	urlChannel := make(chan string)
	go fetchUrlFiles(filename, urlChannel)

	for i := 0; i < thread; i++ {
		waitGroup.Add(1)
		go fetchSongItem(urlChannel)
	}

	waitGroup.Wait()
	opsFinal := atomic.LoadUint64(&count)
	fmt.Println("count:", opsFinal)
}

func fetchSongItem(urlChannel chan string) {
	for {
		select {
		case url := <-urlChannel:
			newCount := atomic.AddUint64(&count, 1)
			if newCount < 6224 {
				continue
			}
			item := strings.Split(url, ",")
			url = item[1]
			songId := item[0]
			oldUrl := url
			var urlType int
			if strings.Contains(url, "hc-audio/hc-audio") {
				url = strings.Replace(url, "hc-audio/hc-audio", "hc-audio", 1)
				urlType = 1
			} else if !strings.Contains(url, "hc-audio") {
				u, _ := URL.Parse(url)
				url = u.Scheme + "://" + u.Host + "/hc-audio" + u.Path
				urlType = 2
			} else {
				urlType = 3
			}
			statusCode := fetchUrl(songId, url)
			if statusCode == 200 {
				fmt.Println(songId, ",", url, ",", statusCode, ",", urlType)
			} else {
				if urlType == 3 {
					oldUrl = strings.Replace(url, "/hc-audio", "", 1)
					statusCode := fetchUrl(songId, oldUrl)
					fmt.Println(songId, ",", oldUrl, ",", statusCode, ",", 4)
				} else {
					statusCode := fetchUrl(songId, oldUrl)
					fmt.Println(songId, ",", oldUrl, ",", statusCode, ",", urlType)
				}

			}
		case <-timeAfter(time.Second * 4):
			waitGroup.Done()
		}
	}
}

func timeAfter(d time.Duration) chan int {
	q := make(chan int, 1)

	time.AfterFunc(d, func() {
		q <- 1
	})

	return q
}

func fetchUrl(songId string, url string) int {
	req, err := http.Head(url)
	if err != nil {
		fmt.Println(songId, ",", url, 4)
		return -1
	}
	defer req.Body.Close()
	return req.StatusCode
}

func fetchUrlFiles(filename string, urlChannel chan string) {
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
