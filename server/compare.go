package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"io"
	"time"
	"sync"
)

type MagicData struct {
	SongId string
	Url string
	StatusCode string
	UrlType string
}

var (
	workerWait = sync.WaitGroup{}
)

func main() {
	var filename string = "t5.log"
	dbUrlChannel := make(chan string)
	var thread int = 5
	go fetchUrlFiles(filename, dbUrlChannel)
	workerWait.Add(1)

	for i := 0; i < thread; i++ {
		go func() {
			for {
				select {
				case dbUrl := <- dbUrlChannel:
					dbItem := parseUrl(dbUrl)
					okItem, err := findSongFromOkSong(dbItem.SongId)
					if err != nil {
						fmt.Println(err)
					}
					compareUrl(okItem.Url, dbItem.Url, okItem.SongId, okItem.StatusCode)
					// fmt.Println(okItem.Url, dbItem.Url)
				case <-timeAfter(time.Second * 4):
					workerWait.Done()
				}
			}
		}()
	}

	workerWait.Wait()
	fmt.Println("done!!!")
}

func compareUrl(okUrl string, dbUrl string, songId string, statusCode string) bool {
	if okUrl != dbUrl {
		fmt.Println(songId, ",", okUrl, ",", dbUrl, ",", statusCode)
		return false
	}
	return true
}

func parseUrl(item string) MagicData {
	items := strings.Split(item, ",")
	statusCode := "0"
	if len(items) > 2 {
		statusCode = items[2]
	}
	return MagicData{
		SongId: strings.Trim(items[0], " "),
		Url: strings.Trim(items[1], " "),
		StatusCode: statusCode,
	}
}

func findSongFromOkSong(songId string) (MagicData, error) {
	fi, err := os.Open("1.txt")
	magicData := MagicData{}
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return magicData, fmt.Errorf("Error: %s", err)
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			return magicData, fmt.Errorf("not found")
			break
		}
		magicData = parseUrl(string(a))
		if magicData.SongId == songId {
			return magicData, nil
		}
	}
	return magicData, nil
}


func timeAfter(d time.Duration) chan int {
	q := make(chan int, 1)

	time.AfterFunc(d, func() {
		q <- 1
	})

	return q
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
