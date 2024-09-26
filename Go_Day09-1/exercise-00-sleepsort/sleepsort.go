package main

import (
	"time"
	"fmt"
)

func sleepSort(nums []int) chan int {
	out := make(chan int)

	// Если входной слайс пуст, сразу закрываем канал
	if len(nums) == 0 {
		close(out)
		return out
	}

	// Запускаем горутины для каждого числа
	for _, num := range nums {
		go func(n int) {
			time.Sleep(time.Duration(n) * time.Second) // Спим количество секунд, равное числу
			out <- n                                  // Отправляем число в канал
		}(num)
	}

	go func() {
		time.Sleep(time.Duration(max(nums)) * time.Second) // Ждём завершения самой длинной паузы
		close(out)                                         // Закрываем канал, когда все горутины завершат работу
	}()

	return out
}

func max(nums []int) int {
	if len(nums) == 0 {
		return 0 // Возвращаем 0, если слайс пуст
	}
	maxVal := nums[0]
	for _, n := range nums {
		if n > maxVal {
			maxVal = n
		}
	}
	return maxVal
}

func main() {
	// Пример входных данных
	input := []int{3, 1, 4, 2}

	// Вызываем sleepSort и читаем результаты
	result := sleepSort(input)

	// Выводим отсортированные значения
	for num := range result {
		fmt.Println(num)
	}
}