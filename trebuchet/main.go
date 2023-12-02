package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	inFlag = flag.String("in", "", "Specify the file to use for input. Defaults to reading from STDIN.")
	debug  = flag.Bool("debug", false, "Turn on debug mode.")
)

func init() {
	flag.BoolFunc("h", "Displays usage", func(_ string) error { flag.CommandLine.Usage(); return nil })
}

func main() {
	flag.Parse()
	if !*debug {
		log.SetOutput(io.Discard)
	}
	log.Println("--- Day 1:", "Trebuchet!? ---")
	f := openFile(*inFlag)
	defer f.Close()
	log.Println("reading input from file:", f.Name())
	fmt.Printf("sum=%d\n", process(f))
}

func openFile(name string) *os.File {
	if name == "" {
		log.Println("expecting input from stdin")
		return os.Stdin
	}
	initial := name
	name = filepath.Clean(name)
	f, err := os.Open(name)
	if err != nil {
		log.Fatalf("error reading input file %q after filepath.Clean %q: %v\n", initial, name, err)
	}
	return f
}

func process(r io.Reader) int {
	var lineNumber, sum int
	s := bufio.NewScanner(r)
	for s.Scan() {
		lineNumber += 1

		line := s.Bytes()
		log.Printf("scanned line #%v: %x %q\n", lineNumber, line, line)
		digits := scanLine(line, filterDigit)
		sum += convert(digits)
	}
	return sum
}

func scanLine(line []byte, filter func(rune) bool) []rune {
	var filtered []rune
	for i, b := range line {
		r := rune(b)
		skip := filter(r)
		log.Printf("scanned rune #%d: %d %q skip? %v\n", i, b, r, skip)
		if skip {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}

func filterDigit(r rune) (skip bool) {
	switch r {
	default:
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return false
	}
}

func convert(digits []rune) int {
	var number string
	for _, r := range digits {
		number += string(r)
	}
	if n := len(number); n > 2 || n == 1 {
		number = string(number[0]) + string(number[n-1])
	}
	v, err := strconv.Atoi(number)
	if err != nil {
		log.Fatalf("error converting ascii to integer: %v\n", err)
	}
	return v
}
