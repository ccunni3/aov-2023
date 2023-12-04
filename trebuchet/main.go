package main

import (
	"flag"
	"log"

	"github.com/ccunni3/aov-2023/trebuchet/sum"
)

var flagInput = flag.String("input", "", "Specify the file to use for input. Defaults to reading from STDIN.")

func main() {
	flag.Parse()
	log.Println("--- Day 1:", "Trebuchet!? ---")
	i, err := sum.File(*flagInput)
	log.Printf("sum=%d error=%v\n", i, err)
}
