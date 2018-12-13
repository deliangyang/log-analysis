package main

import (
	"flag"
	"path/filepath"
	"os/exec"
	"os"
	"fmt"
	"bufio"
	"io"
	"sync"
	"time"
	"strings"
)

func main() {
	var filename string

	flag.StringVar(&filename, "filename", "1.txt", "enter your filename")
	flag.Parse()

	rubyExePath, _ := filepath.Abs("typo\\bin\\ruby.exe")
	urlCrazyPath, _ := filepath.Abs("typo\\urlcrazy")

	waitGroup := sync.WaitGroup{}
	urlList := make(chan string)
	go readUrlByFile(filename, urlList)

	f, _ := os.Create("new-links.txt")
	defer f.Close()

	wr := bufio.NewWriter(f)

	for i := 0; i < 4; i++ {
		waitGroup.Add(1)
		go func() {
			for {
				select {
				case domain := <-urlList:
					println(domain)
					cmd := exec.Command(rubyExePath, urlCrazyPath, domain)
					opBytes, err := cmd.Output()
					if err != nil {
						println(err)
					}
					urls := strings.Split(string(opBytes), "\n")
					extension := getDomainExtension(domain)
					for _, url := range urls {
						println(url)
						if strings.ContainsAny(url, extension) {
							wr.WriteString(url + "\r\n")
						}
					}
					wr.Flush()
				case <-timeoutAfter2(time.Second * 2):
					waitGroup.Done()
				}
			}
		}()
	}
	waitGroup.Wait()
}

func getDomainExtension(domain string) string {
	index := strings.LastIndex(domain, ".")
	return domain[index:]
}

func readUrlByFile(filename string, urlChannel chan string) {
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

func timeoutAfter2(d time.Duration) chan int {
	q := make(chan int, 1)

	time.AfterFunc(d, func() {
		q <- 1
	})

	return q
}
