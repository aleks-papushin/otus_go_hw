package hw03frequencyanalysis

import (
	"sort"
	"unicode/utf8"
)

const topN = 10

var spaces = map[rune]struct{}{' ': {}, '\n': {}, '\t': {}}

func Top10(s string) []string {
	// парсим текст в мапу типа map[string]int (слово:частота)
	quantityMap := parseTextIntoQuantityMap(s)

	// преобразуем quantityMap в мапу типа map[int][]string (n:массив строк со словами, встречающимися n раз)
	// и попутно определяем какова максимальная частота встречаемости
	numStrMap, maxQuantity := projectIntoNumStringMap(quantityMap)

	// начиная с самых часто встречающихся, сортируя массивы лексикографически, взять topN слов из массивов мапы
	outAr := takeTopNSorted(maxQuantity, numStrMap)

	return outAr
}

func parseTextIntoQuantityMap(s string) map[string]int {
	runes := []rune(s)
	words := make(map[string]int)
	textLen := utf8.RuneCountInString(s)

	for i := 0; i < textLen; {
		r := runes[i]
		if _, ok := spaces[r]; ok {
			i++
			continue
		}

		word := []rune{r}
		i++

		for i < textLen {
			r = runes[i]
			if _, ok := spaces[r]; ok {
				i++
				break
			}
			word = append(word, runes[i])
			i++
		}

		wAsString := string(word)
		inMap := words[wAsString] != 0

		if inMap {
			words[wAsString] += 1
		} else {
			words[wAsString] = 1
		}
	}
	return words
}

func projectIntoNumStringMap(quantityMap map[string]int) (map[int][]string, int) {
	numStrMap := make(map[int][]string)
	maxKey := 0 // храним здесь максимальную встреченную частоту
	for k, v := range quantityMap {
		if _, ok := numStrMap[v]; ok {
			slc := numStrMap[v]
			slc = append(slc, k)
			numStrMap[v] = slc
		} else {
			numStrMap[v] = []string{k}
			if v > maxKey {
				maxKey = v
			}
		}
	}
	return numStrMap, maxKey
}

func takeTopNSorted(maxKey int, numStrMap map[int][]string) []string {
	outAr := []string{}
	sum := 0
	for k := maxKey; k > 0 && sum < topN; k-- {
		if ar, ok := numStrMap[k]; ok {
			sort.Strings(ar)
			for _, w := range ar {
				outAr = append(outAr, w)
				sum++
				if sum >= topN {
					break
				}
			}
		}
	}
	return outAr
}
