package main

import (
	"os"
	"bufio"
	"io"
	"strings"
	"strconv"
	"fmt"
)

type MagicData struct {
	SongId string
	Url string
	StatusCode string
	UrlType string
}

func main() {
	filename := "1.txt"
	fmt.Println()
	fetchFile(filename)
}

func fetchFile(filename string)  {
	f, _ := os.Open(filename)
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		item := strings.Split(string(a), ",")
		if len(item) < 4 {
			// fmt.Println(string(a))
			continue
		}
		magicData := MagicData{
			SongId: item[0],
			Url: item[1],
			StatusCode: item[2],
			UrlType: item[3],
		}
		printAvailableSongUrl(magicData)
	}
}

func printAvailableSongUrl(data MagicData) {
	urlType, _ := strconv.Atoi(strings.Trim(data.UrlType, " "))
	if urlType == 2 && strings.Contains(data.Url, "hc-audio") {
		fmt.Print(data.SongId, ", ")
		return
	} else if urlType == 4 {

	} else if urlType == 3 {
		return
	} else if urlType == 1 {

	}
}
