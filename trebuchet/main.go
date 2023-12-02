package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var (
	inputFilepath = flag.String("in", "", "Specify the file to use for input. Defaults to reading from STDIN.")
)

func init() {
	flag.BoolFunc("h", "Displays usage", func(_ string) error { flag.CommandLine.Usage(); return nil })
}

func main() {
	flag.Parse()
	log.Println("--- Day 1:", "Trebuchet!? ---")
	validateInputFilePath(inputFilepath)
	log.Println("reading input from file:", *inputFilepath)
}

func validateInputFilePath(s *string) {
	if *s == "" {
		log.Println("reading input from stdin")
		return
	}
	p := filepath.Clean(*s)
	info, err := os.Stat(p)
	if err != nil {
		log.Fatalf("error reading input file %q: %v\n", *inputFilepath, err)
	}
	if info.IsDir() {
		log.Fatalln("input file path cannot be a directory")
	}
	s = &p
}
