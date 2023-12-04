package sum

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"unicode/utf8"
)

// File parses the named file and computes the sum.
func File(name string) (int64, error) {
	var (
		file *os.File
		err  error
	)
	if name == "" {
		file = os.Stdin
	} else {
		name = filepath.Clean(name)
		file, err = os.Open(name)
		if err != nil {
			return 0, err
		}
		defer file.Close()
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(defaultScanDigits)

	s := &Sum{
		scanner: scanner,
	}

	return s.Sum()
}

func New(r io.Reader, split bufio.SplitFunc, parse ParseFunc) *Sum {
	s := bufio.NewScanner(r)
	if split == nil {
		split = defaultScanDigits
	}
	s.Split(split)
	return &Sum{
		scanner:   s,
		parseFunc: parse,
	}
}

type ParseFunc func(string) (int64, error)

type Sum struct {
	scanner   *bufio.Scanner
	parseFunc ParseFunc
}

func (s *Sum) Sum() (int64, error) {
	var sum int64
	for s.scanner.Scan() {
		num := s.scanner.Text()
		i, err := s.parse(num)
		if err != nil {
			return 0, err
		}
		sum += i
	}
	return sum, nil
}

func (s *Sum) parse(num string) (int64, error) {
	if s.parseFunc == nil {
		i, err := strconv.Atoi(num)
		return int64(i), err
	}
	return s.parseFunc(num)
}

var defaultScanDigits = ScanFirstAndLastDigitLiterals

// ScanFirstAndLastDigitLiterals is a split function for a bufio.Scanner that returns a
// two-digit number for each line of text. The number is computed by concatenating the
// first and last digit literals found on the line. If only one digit is found, it is doubled.
func ScanFirstAndLastDigitLiterals(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Split data into a line.
	lineLen, line, err := bufio.ScanLines(data, atEOF)
	if lineLen == 0 {
		return lineLen, line, err
	}

	// Find each digit literal in the line.
	for i, width := 0, 0; i < len(line); i += width {
		var r rune
		r, width = utf8.DecodeRune(line[i:])
		if r >= '0' && r <= '9' {
			token = utf8.AppendRune(token, r)
		}
	}

	token = concatFirstAndLastTokens(token)
	if token == nil {
		// Request more data.
		return 0, nil, nil
	}
	return lineLen, token, nil
}

func concatFirstAndLastTokens(tokens []byte) []byte {
	// Compute the number from the literals.
	tokenCount := utf8.RuneCount(tokens)
	if tokenCount == 0 {
		// Request more data.
		return nil
	}

	first, firstWidth := utf8.DecodeRune(tokens)
	if tokenCount == 1 {
		tokens = utf8.AppendRune(tokens, first)
		return tokens
	}

	// Take the first and last digits and concatenate them.
	_, lastWidth := utf8.DecodeLastRune(tokens)
	tokens = append(tokens[:firstWidth], tokens[len(tokens)-lastWidth:]...)
	return tokens
}
