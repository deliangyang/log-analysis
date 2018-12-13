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
	"container/heap"
)

type Heap []string

//构造的是小顶堆，大顶堆只需要改一下下面的符号
func (h *Heap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *Heap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *Heap) Len() int {
	return len(*h)
}

func (h *Heap) Pop() interface{} {
	x := (*h)[h.Len() - 1]
	*h = (*h)[: h.Len() - 1]
	return x
}

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(string))
}

func (h *Heap) Remove(idx int) interface{} {
	h.Swap(idx, h.Len() - 1)
	return h.Pop()
}

var (
	pattern, _ = regexp.Compile(`反向链接：</span>(\d+)</li>`)
	h = &Heap{}
)

func main() {
	var filename string
	var sleepTime int64

	heap.Init(h)

	flag.StringVar(&filename, "filename", "", "enter your filename")
	flag.Int64Var(&sleepTime, "sleepTime", 3, "fetch domain sleep time")
	flag.Parse()

	if filename == "" {
		os.Exit(-1)
	}

	waitGroup := sync.WaitGroup{}
	domainList := make(chan string)
	go readFile(filename, domainList)

	f, _ := os.Create("links.txt")
	defer f.Close()

	wr := bufio.NewWriter(f)


	for i := 0; i < 1; i++ {
		waitGroup.Add(1)
		go func(domainList chan string, sleepTime int64) {
			for {
				select {
				case domain := <-domainList:
					lines, isFetch := fetchPage(domain, sleepTime)
					if isFetch {
						fmt.Println(domain + "\t" + lines)
						wr.WriteString(domain + "\t" + lines + "\r\n")
						wr.Flush()
					} else {
						fmt.Println(domain + "\t" + "未获取到链接数，重新抓取")
					}
				case <-timeoutAfter(time.Second * 10):
					waitGroup.Done()
				}
			}
		}(domainList, sleepTime)
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

func fetchPage(url string, sleepTime int64) (links string, isFetch bool) {
	targetUrl := "http://alexa.chinaz.com/?domain=" + url
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return "", false
	}
	client := &http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	value := pattern.FindSubmatch(body)
	linkNum := "未获取到, 重新抓取"
	if len(value) >= 2 {
		linkNum = string(value[1])
	} else {
		h.Push(url)
		return "0", false
	}
	time.Sleep(time.Second * time.Duration(sleepTime))
	return linkNum, true
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

	for h.Len() > 0 {
		urlChannel <- heap.Pop(h).(string)
	}
}
