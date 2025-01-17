package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var months = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

var (
	flagColumn         = flag.Int("k", 0, "Указать номер колонки (1-based), по которой сортировать (0 = без использования колонки)")
	flagNumeric        = flag.Bool("n", false, "Сортировать по числовому значению")
	flagReverse        = flag.Bool("r", false, "Сортировать в обратном порядке")
	flagUnique         = flag.Bool("u", false, "Не выводить повторяющиеся строки")
	flagMonth          = flag.Bool("M", false, "Сортировать по названию месяца (jan, feb, mar, ...)")
	flagIgnoreTrailing = flag.Bool("b", false, "Игнорировать хвостовые пробелы")
	flagCheckSorted    = flag.Bool("c", false, "Проверять, отсортированы ли данные (без вывода)")
	flagHuman          = flag.Bool("h", false, "Сортировать с учётом суффиксов (K, M и т.д.)")
)

func getKey(line string, col int) string {
	if col <= 0 {
		return line
	}
	fields := strings.Fields(line)
	if len(fields) < col {
		return ""
	}
	return fields[col-1]
}

func parseHumanSize(s string) (float64, error) {
	s = strings.ToUpper(s)
	multiplier := 1.0
	switch {
	case strings.HasSuffix(s, "K"):
		multiplier = 1000
		s = strings.TrimSuffix(s, "K")
	case strings.HasSuffix(s, "M"):
		multiplier = 1000_000
		s = strings.TrimSuffix(s, "M")
	case strings.HasSuffix(s, "G"):
		multiplier = 1000_000_000
		s = strings.TrimSuffix(s, "G")
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return val * multiplier, nil
}

func compareStrings(a, b string) int {
	ka := getKey(a, *flagColumn)
	kb := getKey(b, *flagColumn)

	if *flagIgnoreTrailing {
		ka = strings.TrimRight(ka, " \t")
		kb = strings.TrimRight(kb, " \t")
	}

	if *flagMonth {
		ma, okA := months[strings.ToLower(ka)]
		mb, okB := months[strings.ToLower(kb)]
		if okA && okB {
			switch {
			case ma < mb:
				return -1
			case ma > mb:
				return 1
			default:
				return 0
			}
		}
		return strings.Compare(ka, kb)
	}

	if *flagNumeric {
		fa, errA := strconv.ParseFloat(ka, 64)
		fb, errB := strconv.ParseFloat(kb, 64)
		if errA == nil && errB == nil {
			switch {
			case fa < fb:
				return -1
			case fa > fb:
				return 1
			default:
				return 0
			}
		}
		return strings.Compare(ka, kb)
	}

	if *flagHuman {
		ha, errA := parseHumanSize(ka)
		hb, errB := parseHumanSize(kb)
		if errA == nil && errB == nil {
			switch {
			case ha < hb:
				return -1
			case ha > hb:
				return 1
			default:
				return 0
			}
		}
		return strings.Compare(ka, kb)
	}

	return strings.Compare(ka, kb)
}

func isSorted(lines []string) bool {
	for i := 1; i < len(lines); i++ {
		if compareStrings(lines[i-1], lines[i]) > 0 {
			return false
		}
	}
	return true
}

func main() {
	flag.Parse()

	var in io.Reader = os.Stdin
	if flag.NArg() > 0 {
		fileName := flag.Arg(0)
		f, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("Не удалось открыть файл: %v", err)
		}
		defer f.Close()
		in = f
	} else {
		log.Println("Чтение из stdin (по умолчанию)")
	}

	scanner := bufio.NewScanner(in)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Ошибка чтения: %v", err)
	}

	if *flagCheckSorted {
		if isSorted(lines) {
			fmt.Println("Данные уже отсортированы")
			os.Exit(0)
		} else {
			fmt.Println("Данные НЕ отсортированы")
			os.Exit(1)
		}
	}

	sort.SliceStable(lines, func(i, j int) bool {
		return compareStrings(lines[i], lines[j]) < 0
	})

	if *flagReverse {
		for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
			lines[i], lines[j] = lines[j], lines[i]
		}
	}

	if *flagUnique && len(lines) > 0 {
		uniq := lines[:1]
		for i := 1; i < len(lines); i++ {
			if compareStrings(lines[i], uniq[len(uniq)-1]) != 0 {
				uniq = append(uniq, lines[i])
			}
		}
		lines = uniq
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}
