// Package moneybag предоставляет функции для работы с минимальным количеством монет.
//
// Основное отличие функции minCoins2 от оригинальной minCoins заключается в том, что:
// - minCoins2 обрабатывает списки монет с дубликатами и неотсортированными номиналами;
// - она также корректно работает с пустым списком номиналов, возвращая пустой результат.
//
// Оптимизация:
// - Используется сортировка списка номиналов, если он не отсортирован.
// - Применена проверка дубликатов для уменьшения объема вычислений.
package main

import (
	"sort"
)

// minCoins2 находит минимальное количество монет для заданной суммы.
// Она исправляет недостатки оригинальной функции minCoins:
//  - Обрабатывает неотсортированные и дублированные номиналы.
//  - Возвращает пустой массив, если список номиналов пуст.
// Порядок аргументов:
// - val: сумма, которую нужно составить.
// - coins: список номиналов монет (в любом порядке).
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
