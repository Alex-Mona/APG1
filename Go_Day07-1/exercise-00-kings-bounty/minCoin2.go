// minCoins2.go
package main

import (
	"sort"
)

// minCoins2 исправляет недостатки оригинальной функции minCoins.
// Функция работает корректно с несортированными и дублированными номиналами.
// Если список номиналов пуст, функция возвращает пустой список.
func minCoins2(val int, coins []int) []int {
	if len(coins) == 0 {
		return []int{}
	}

	// Удаление дубликатов и сортировка номиналов
	uniqueCoins := removeDuplicatesAndSort(coins)

	res := make([]int, 0)
	for i := len(uniqueCoins) - 1; i >= 0; i-- {
		for val >= uniqueCoins[i] {
			val -= uniqueCoins[i]
			res = append(res, uniqueCoins[i])
		}
	}

	return res
}

// Удаление дубликатов и сортировка списка номиналов
func removeDuplicatesAndSort(coins []int) []int {
	uniqueMap := make(map[int]bool)
	for _, coin := range coins {
		uniqueMap[coin] = true
	}

	uniqueCoins := make([]int, 0, len(uniqueMap))
	for coin := range uniqueMap {
		uniqueCoins = append(uniqueCoins, coin)
	}

	sort.Ints(uniqueCoins)
	return uniqueCoins
}
