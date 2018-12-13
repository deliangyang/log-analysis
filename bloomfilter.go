package main

import (
	"github.com/willf/bitset"
	"flag"
	"os"
	"bufio"
	"io"
	"fmt"
)

const DefaultSize = 2 << 24

var seeds = []uint{7, 11, 13, 31, 37, 61}

type BloomFilter struct {
	set   *bitset.BitSet
	funcs [6]SimpleHash
}

func NewBloomFilter() *BloomFilter {
	bf := new(BloomFilter)
	for i := 0; i < len(bf.funcs); i++ {
		bf.funcs[i] = SimpleHash{DefaultSize, seeds[i]}
	}
	bf.set = bitset.New(DefaultSize)
	return bf
}

func (bf BloomFilter) add(value string) {
	for _, f := range bf.funcs {
		bf.set.Set(f.hash(value))
	}
}

func (bf BloomFilter) contains(value string) bool {
	if value == "" {
		return false
	}
	ret := true
	for _, f := range bf.funcs {
		ret = ret && bf.set.Test(f.hash(value))
	}
	return ret
}

type SimpleHash struct {
	cap  uint
	seed uint
}

func (s SimpleHash) hash(value string) uint {
	var result uint = 0
	for i := 0; i < len(value); i++ {
		result = result*s.seed + uint(value[i])
	}
	return (s.cap - 1) & result
}

func main() {
	var filename string
	flag.StringVar(&filename,  "filename", "", "enter filename")
	flag.Parse()

	if filename == "" {
		os.Exit(-1)
	}

	filter := NewBloomFilter()
	f, _ := os.Open(filename)
	defer f.Close()

	wf, _ := os.Create("1.txt")
	defer wf.Close()
	bw := bufio.NewWriter(wf)

	br := bufio.NewReader(f)
	count := 0
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		number := string(a)
		if filter.contains(number) {
			continue
		}
		count++
		if count % 100 == 2 {
			bw.Flush()
			fmt.Println("contiue run")
		}
		filter.add(number)
		fmt.Println(number)
		bw.WriteString(number + "\n")
	}
	bw.Flush()
}