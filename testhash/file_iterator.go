package main


import "fmt"

func main() {
	hash := make(map[string]string)
	hash["aaa"] = "hello world"

	fmt.Println(hash["a"] == "")
}