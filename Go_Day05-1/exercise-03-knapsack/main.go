package main

import "fmt"

// Определение структуры Present
type Present struct {
    Value int
    Size  int
}

// Функция для решения задачи о рюкзаке
func grabPresents(presents []Present, capacity int) []Present {
    n := len(presents)
    dp := make([][]int, n+1)
    for i := range dp {
        dp[i] = make([]int, capacity+1)
    }

    for i := 1; i <= n; i++ {
        for w := 0; w <= capacity; w++ {
            if presents[i-1].Size <= w {
                dp[i][w] = max(dp[i-1][w], dp[i-1][w-presents[i-1].Size]+presents[i-1].Value)
            } else {
                dp[i][w] = dp[i-1][w]
            }
        }
    }

    // Обратный ход для определения списка подарков
    result := []Present{}
    w := capacity
    for i := n; i > 0 && w > 0; i-- {
        if dp[i][w] != dp[i-1][w] {
            result = append(result, presents[i-1])
            w -= presents[i-1].Size
        }
    }

    return result
}

// Вспомогательная функция для нахождения максимума из двух чисел
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func main() {
    // Пример данных
    presents := []Present{
        {Value: 10, Size: 5},
        {Value: 40, Size: 4},
        {Value: 30, Size: 6},
        {Value: 50, Size: 3},
    }
    capacity := 10

    // Вызов функции grabPresents
    selectedPresents := grabPresents(presents, capacity)

    // Вывод результата
    fmt.Println("Выбранные подарки:")
    for _, present := range selectedPresents {
        fmt.Printf("Value: %d, Size: %d\n", present.Value, present.Size)
    }
}
