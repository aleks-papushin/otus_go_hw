package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const topN = 10

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
	wordsArray := strings.Fields(s)
	quantityMap := make(map[string]int)

	for _, w := range wordsArray {
		quantityMap[w] += 1
	}

	return quantityMap
}

func projectIntoNumStringMap(quantityMap map[string]int) (map[int][]string, int) {
	numStrMap := make(map[int][]string)
	maxQ := 0 // храним здесь максимальную встреченную частоту
	for k, v := range quantityMap {
		if _, ok := numStrMap[v]; ok {
			slc := numStrMap[v]
			slc = append(slc, k)
			numStrMap[v] = slc
		} else {
			numStrMap[v] = []string{k}
			if v > maxQ {
				maxQ = v
			}
		}
	}
	return numStrMap, maxQ
}

func takeTopNSorted(maxQ int, numStrMap map[int][]string) []string {
	var outSlc []string
	sum := 0
	for k := maxQ; k > 0 && sum < topN; k-- {
		if ar, ok := numStrMap[k]; ok {
			sort.Strings(ar)
			for _, w := range ar {
				outSlc = append(outSlc, w)
				sum++
				if sum >= topN {
					break
				}
			}
		}
	}
	return outSlc
}
