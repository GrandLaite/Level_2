/*
Написать функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' — принадлежат одному множеству;
'листок', 'слиток' и 'столик' — другому.

Требования
1.Входные данные для функции: ссылка на массив,
каждый элемент которого — слово на русском языке в кодировке utf8.
2.Выходные данные: ссылка на мапу множеств анаграмм.
3.Ключ — первое встретившееся в словаре слово из множества. Значение — ссылка на массив,
каждый элемент которого, слово из множества.
4.Массив должен быть отсортирован по возрастанию.
5.Множества из одного элемента не должны попасть в результат.
6.Все слова должны быть приведены к нижнему регистру.
7.В результате каждое слово должно встречаться только один раз.
*/

package main

import (
	"fmt"
	"sort"
	"strings"
)

func FindAnagrams(words []string) map[string][]string {

	type groupInfo struct {
		firstSeen   string
		uniqueWords map[string]struct{}
	}

	groups := make(map[string]*groupInfo)

	for _, originalWord := range words {
		w := strings.ToLower(strings.TrimSpace(originalWord))
		if w == "" {
			continue
		}
		sorted := sortRunes(w)

		if _, ok := groups[sorted]; !ok {
			groups[sorted] = &groupInfo{
				firstSeen:   w,
				uniqueWords: make(map[string]struct{}),
			}
		}

		groups[sorted].uniqueWords[w] = struct{}{}
	}

	result := make(map[string][]string)
	for _, info := range groups {
		var arr []string
		for w := range info.uniqueWords {
			arr = append(arr, w)
		}
		sort.Strings(arr)

		if len(arr) < 2 {
			continue
		}

		result[info.firstSeen] = arr
	}

	return result
}

func sortRunes(s string) string {
	r := []rune(s)
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return string(r)
}

func main() {
	words := []string{
		"пятак", "пятка", "тяпка",
		"листок", "Слиток", "столик",
		"кот", "ток", "окт",
	}

	anagrams := FindAnagrams(words)
	for key, group := range anagrams {
		fmt.Printf("%q: %v\n", key, group)
	}
}
