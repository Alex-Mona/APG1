package main

import (
    "container/heap"
    "fmt"
)

// Определение структуры Present
type Present struct {
    Value int
    Size  int
}

// PresentHeap реализует интерфейс heap.Interface
type PresentHeap []Present

func (h PresentHeap) Len() int { return len(h) }
func (h PresentHeap) Less(i, j int) bool {
    if h[i].Value == h[j].Value {
        return h[i].Size < h[j].Size
    }
    return h[i].Value > h[j].Value
}
func (h PresentHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *PresentHeap) Push(x interface{}) {
    *h = append(*h, x.(Present))
}

func (h *PresentHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

// Функция для получения n самых крутых подарков
func getNCoolestPresents(presents []Present, n int) ([]Present, error) {
    if n > len(presents) || n < 0 {
        return nil, fmt.Errorf("некорректное значение n")
    }

    h := &PresentHeap{}
    heap.Init(h)

    for _, p := range presents {
        heap.Push(h, p)
    }

    coolest := make([]Present, n)
    for i := 0; i < n; i++ {
        coolest[i] = heap.Pop(h).(Present)
    }

    return coolest, nil
}

func main() {
    // Пример использования
    presents := []Present{
        {Value: 5, Size: 1},
        {Value: 4, Size: 5},
        {Value: 3, Size: 1},
        {Value: 5, Size: 2},
    }

    coolest, err := getNCoolestPresents(presents, 2)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(coolest) // Вывод: [{5 1} {5 2}]
    }
}
