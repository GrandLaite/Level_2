/*
Реализовать утилиту фильтрации по аналогии с консольной утилитой
(man grep — смотрим описание и основные параметры).

Реализовать поддержку утилитой следующих ключей:
-A - "after": печатать +N строк после совпадения;
-B - "before": печатать +N строк до совпадения;
-C - "context": (A+B) печатать ±N строк вокруг совпадения;
-c - "count": количество строк;
-i - "ignore-case": игнорировать регистр;
-v - "invert": вместо совпадения, исключать;
-F - "fixed": точное совпадение со строкой, не паттерн;
-n - "line num": напечатать номер строки.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type grepOptions struct {
	after      int  // -A
	before     int  // -B
	countOnly  bool // -c
	ignoreCase bool // -i
	invert     bool // -v
	fixed      bool // -F
	lineNum    bool // -n
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: mygrep [options] PATTERN [FILE]

Options:
  -A int     print +N lines after match
  -B int     print +N lines before match
  -C int     print ±N lines around match (shorthand for -A N -B N)
  -c         print only count of matching lines
  -i         ignore case
  -v         invert match (select non-matching lines)
  -F         fixed string match (not substring)
  -n         print line number
`)
	os.Exit(1)
}

func main() {
	var opts grepOptions
	var context int

	flag.IntVar(&opts.after, "A", 0, "print +N lines after match")
	flag.IntVar(&opts.before, "B", 0, "print +N lines before match")
	flag.IntVar(&context, "C", 0, "print ±N lines around match (shorthand for -A N -B N)")
	flag.BoolVar(&opts.countOnly, "c", false, "print only count of matching lines")
	flag.BoolVar(&opts.ignoreCase, "i", false, "ignore case")
	flag.BoolVar(&opts.invert, "v", false, "invert match (select non-matching lines)")
	flag.BoolVar(&opts.fixed, "F", false, "fixed string match")
	flag.BoolVar(&opts.lineNum, "n", false, "print line number")
	flag.Parse()

	if context > 0 {
		opts.after = context
		opts.before = context
	}

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
	}
	pattern := args[0]

	var in *os.File
	if len(args) > 1 {
		f, err := os.Open(args[1])
		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}
		defer f.Close()
		in = f
	} else {
		in = os.Stdin
	}

	scanner := bufio.NewScanner(in)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	matchedIdx := matchLines(lines, pattern, opts)

	if opts.countOnly {
		fmt.Println(len(matchedIdx))
		return
	}

	printMatched(lines, matchedIdx, opts)
}

func matchLines(lines []string, pattern string, opts grepOptions) []int {
	var matchedIdx []int

	if opts.ignoreCase {
		pattern = strings.ToLower(pattern)
	}

	for i, line := range lines {
		lineToCheck := line
		if opts.ignoreCase {
			lineToCheck = strings.ToLower(line)
		}

		var isMatch bool
		if opts.fixed {
			isMatch = (lineToCheck == pattern)
		} else {
			isMatch = strings.Contains(lineToCheck, pattern)
		}

		if opts.invert {
			isMatch = !isMatch
		}

		if isMatch {
			matchedIdx = append(matchedIdx, i)
		}
	}
	return matchedIdx
}

func printMatched(lines []string, matchedIdx []int, opts grepOptions) {
	matchesSet := make(map[int]bool, len(matchedIdx))
	for _, idx := range matchedIdx {
		matchesSet[idx] = true
	}

	lastPrinted := -1

	for _, idx := range matchedIdx {
		start := idx - opts.before
		if start < 0 {
			start = 0
		}
		end := idx + opts.after
		if end >= len(lines)-1 {
			end = len(lines) - 1
		}

		for i := start; i <= end; i++ {
			if i <= lastPrinted {
				continue
			}
			if opts.lineNum {
				fmt.Printf("%d:%s\n", i+1, lines[i])
			} else {
				fmt.Println(lines[i])
			}
			lastPrinted = i
		}
	}
}
