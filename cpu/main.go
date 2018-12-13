package main

import (
	"os"
	"fmt"
	"flag"
	"os/exec"
	"net/http"
	"io/ioutil"
	"time"
	"github.com/aurelien-rainone/assertgo"
	"sync"
)

func test() int {
	println("2")
	return 4
}

func main() {
	assert.True(1 == 1)
	ch := make(chan int)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		time.Sleep(time.Second * 2)
		v := <-ch
		fmt.Println("v:", v)
		waitGroup.Done()
	}()
	ch <- 1
	waitGroup.Wait()
	fmt.Println(2)
	os.Exit(1)
	req, err := http.Get("http://www.baidu.com")
	if err != nil {
		fmt.Println(err)
	}

	rsq, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(rsq))


	// err := errors.New("executable file not found in %PATH%")
	fmt.Println(err)
	dataCmd := exec.Command("date")
	output, err := dataCmd.Output()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
	os.Exit(1)

	argWithProg := os.Args
	argWithoutProg := os.Args[1:]

	wordPtr := flag.String("word", "foo", "a String")
	numberPtr := flag.Int("number", 1, "an int")

	var test string
	flag.StringVar(&test, "test", "tt", "test string")
	flag.Parse()
	fmt.Println(*wordPtr)
	fmt.Println(*numberPtr)
	fmt.Println(test)

	arg := os.Args[3]

	fmt.Println(argWithProg)
	fmt.Println(argWithoutProg)
	fmt.Println(arg)
}
