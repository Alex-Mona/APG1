package main

import (
    "errors"
    "fmt"
	"unsafe"
)

func getElement(arr []int, idx int) (int, error) {
    // Проверяем на валидность слайс и индекс
    if len(arr) == 0 {
        return 0, errors.New("слайс пуст")
    }
    if idx < 0 {
        return 0, errors.New("индекс не может быть отрицательным")
    }
    if idx >= len(arr) {
        return 0, errors.New("индекс выходит за границы слайса")
    }
    
    // Возвращаем элемент через указатель на первый элемент слайса
    return *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&arr[0])) + uintptr(idx)*unsafe.Sizeof(arr[0]))), nil
}

func main() {
    arr := []int{10, 20, 30, 40, 50}
    elem, err := getElement(arr, 3)
    if err != nil {
        fmt.Println("Ошибка:", err)
    } else {
        fmt.Println("Элемент:", elem)
    }
}
