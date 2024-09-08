package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Структуры данных, которые будут использоваться для обработки запросов и формирования ответов CandyOrder, ResponseSuccess, ResponseError
type CandyOrder struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

type ResponseSuccess struct {
	Thanks string `json:"thanks"`
	Change int    `json:"change"`
}

type ResponseError struct {
	Error string `json:"error"`
}

// Цены на конфеты
var candyPrices = map[string]int{
	"CE": 10,
	"AA": 15,
	"NT": 20,
	"DE": 25,
	"YR": 30,
}

// Функция, которая будет обрабатывать POST-запросы на путь /buy_candy
func buyCandy(w http.ResponseWriter, r *http.Request) {
	var order CandyOrder

	// Декодируем JSON-запрос
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверка на валидность типа конфет
	price, exists := candyPrices[order.CandyType]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{Error: "Invalid candy type"})
		return
	}

	// Проверка на валидность количества конфет
	if order.CandyCount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{Error: "Candy count must be positive"})
		return
	}

	// Расчет необходимой суммы
	totalCost := price * order.CandyCount

	// Проверка на достаточность денег
	if order.Money < totalCost {
		w.WriteHeader(http.StatusPaymentRequired)
		json.NewEncoder(w).Encode(ResponseError{Error: fmt.Sprintf("You need %d more money!", totalCost-order.Money)})
		return
	}

	// Если все ок, возвращаем успешный ответ
	change := order.Money - totalCost
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ResponseSuccess{Thanks: "Thank you!", Change: change})
}

// Настройка маршрута и запуск HTTP сервера
func main() {
	http.HandleFunc("/buy_candy", buyCandy)
	fmt.Println("Server is running on http://127.0.0.1:3333")
	http.ListenAndServe(":3333", nil)
}
