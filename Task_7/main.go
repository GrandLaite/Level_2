/*
Реализовать утилиту аналог консольной команды cut (man cut).
Утилита должна принимать строки через STDIN,
разбивать по разделителю (TAB) на колонки и выводить запрошенные.

Реализовать поддержку утилитой следующих ключей:
-f — "fields": выбрать поля (колонки);
-d — "delimiter": использовать другой разделитель;
-s — "separated": только строки с разделителем.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type cutOptions struct {
	fields    string // "1,2,5" и т.п.
	delimiter string
	separated bool // -s
}

func parseFields(fieldsStr string) ([]int, error) {
	if fieldsStr == "" {
		return nil, fmt.Errorf("no fields specified")
	}
	parts := strings.Split(fieldsStr, ",")
	var fields []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		val, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid field number: %q", p)
		}
		if val < 1 {
			return nil, fmt.Errorf("field number must be >= 1: %d", val)
		}
		fields = append(fields, val-1)
	}
	return fields, nil
}

func main() {
	var opts cutOptions

	flag.StringVar(&opts.fields, "f", "", "выбрать поля (колонки), например '1,2,5'")
	flag.StringVar(&opts.delimiter, "d", "\t", "использовать другой разделитель (по умолчанию табуляция)")
	flag.BoolVar(&opts.separated, "s", false, "не выводить строки, не содержащие разделитель")

	flag.Parse()

	fieldIndexes, err := parseFields(opts.fields)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка в флагах:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		columns := strings.Split(line, opts.delimiter)
		if opts.separated && len(columns) < 2 {
			continue
		}

		var selected []string
		for _, colIndex := range fieldIndexes {
			if colIndex < len(columns) {
				selected = append(selected, columns[colIndex])
			} else {
			}
		}

		fmt.Println(strings.Join(selected, opts.delimiter))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения stdin:", err)
		os.Exit(1)
	}
}
