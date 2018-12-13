package main

import (
	"flag"
	"net/http"
	"io/ioutil"
	"io"
	"fmt"
)

type Downloader struct {
	Url string
	UrlMd5 string
	Type uint8
}

func (downloader Downloader) CheckHasDownload() {

}

func (downloader Downloader) DownloadFile() bool {
	req, err := http.Get(downloader.Url)
	if err != nil {
		return false
	}
	defer req.Body.Close()
	_, err = io.Copy(newFile, req.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
}

var (
	filename string
	thread int
	line int
)

func main() {
	flag.IntVar(&thread, "thread", 3, "number of worker")
	flag.IntVar(&line, "line", 0, "the start line of excel")
	flag.StringVar(&filename, "filename", "", "the excel file name")
	flag.Parse()

	if len(filename) <= 0 {
		panic("need excel file name")
	}
}
