package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
)

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Функция C для создания ASCII-арт изображения с фразой
char *ask_cow(const char *phrase) {
    int phrase_len = strlen(phrase);
    // Выделяем память для буфера с учетом длины фразы и отступов
    char *buf = (char *)malloc(sizeof(char) * (160 + (phrase_len + 2) * 3));
    if (!buf) return NULL; // Проверка успешного выделения памяти
    strcpy(buf, " ");

    // Добавляем верхнюю границу
    for (unsigned int i = 0; i < phrase_len + 2; ++i) {
        strcat(buf, "_");
    }

    strcat(buf, "\n< ");
    strcat(buf, phrase);
    strcat(buf, " ");
    strcat(buf, ">\n ");

    // Добавляем нижнюю границу
    for (unsigned int i = 0; i < phrase_len + 2; ++i) {
        strcat(buf, "-");
    }
    strcat(buf, "\n");
    strcat(buf, "        \\   ^__^\n");
    strcat(buf, "         \\  (oo)\\_______\n");
    strcat(buf, "            (__)\\       )\\/\\\n");
    strcat(buf, "                ||----w |\n");
    strcat(buf, "                ||     ||\n");

    return buf;
}
*/
import "C"
import (
	"bytes"
	"unsafe"
)

// askCow вызывает функцию C для создания ASCII-арт изображения с фразой
func askCow(phrase string) string {
	// Преобразуем Go строку в C строку
	cPhrase := C.CString(phrase)
	defer C.free(unsafe.Pointer(cPhrase)) // Освобождаем память после использования

	// Вызываем функцию C
	ptr := C.ask_cow(cPhrase)
	if ptr == nil {
		return "Error in C function" // Обработка ошибки, если функция C вернула NULL
	}
	defer C.free(unsafe.Pointer(ptr)) // Освобождаем память, выделенную функцией C

	// Преобразуем C буфер в Go строку
	b := C.GoBytes(unsafe.Pointer(ptr), C.int(160+len(phrase)+4)) // Корректировка длины, если нужно
	b = bytes.TrimRight(b, "\x00")                                // Удаляем нулевые байты

	return string(b)
}

// Data представляет структуру данных для запроса
type Data struct {
	Money      int    `json:"money"`      // Количество денег
	CandyType  string `json:"candyType"`  // Тип конфет
	CandyCount int    `json:"candyCount"` // Количество конфет
}

// Функция для получения цены конфет по их типу
func getPrice(candyType string) (int, error) {
	switch candyType {
	case "CE":
		return 10, nil
	case "AA":
		return 15, nil
	case "NT":
		return 17, nil
	case "DE":
		return 21, nil
	case "YR":
		return 23, nil
	default:
		return 0, errors.New("wrong candy type") // Возвращаем ошибку, если тип конфет неверный
	}
}

// Обработчик HTTP-запросов
func handler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса - POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed) // Возвращаем ошибку 405, если метод не поддерживается
		return
	}

	var data Data
	// Декодируем JSON-данные из тела запроса
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding request body:", err)             // Логируем ошибку декодирования
		http.Error(w, "Invalid request body", http.StatusBadRequest) // Возвращаем ошибку 400, если тело запроса некорректно
		return
	}

	var response interface{}
	candyPrice, err := getPrice(data.CandyType) // Получаем цену конфет

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Возвращаем ошибку 400, если тип конфет неверный
		response = struct {
			Error string `json:"error"`
		}{"Wrong candy type!"}
	} else if data.CandyCount < 0 {
		w.WriteHeader(http.StatusBadRequest) // Возвращаем ошибку 400, если количество конфет отрицательное
		response = struct {
			Error string `json:"error"`
		}{"Negative candy count!"}
	} else if candyPrice*data.CandyCount > data.Money {
		amount := candyPrice*data.CandyCount - data.Money
		w.WriteHeader(http.StatusPaymentRequired) // Возвращаем ошибку 402, если недостаточно денег
		response = struct {
			Error string `json:"error"`
		}{"You need " + strconv.Itoa(amount) + " more money!"}
	} else {
		change := data.Money - candyPrice*data.CandyCount
		thanks := askCow("Thank you!") // Используем C-функцию для генерации сообщения
		w.WriteHeader(http.StatusOK)   // Возвращаем успешный ответ 200
		response = struct {
			Change int    `json:"change"`
			Thanks string `json:"thanks"`
		}{change, thanks}
	}

	w.Header().Set("Content-Type", "application/json") // Устанавливаем тип содержимого ответа
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding response:", err) // Логируем ошибку кодирования ответа
	}
}

// Основная функция для запуска сервера
func main() {
	http.HandleFunc("/buy_candy", handler) // Регистрируем обработчик для маршрута /buy_candy

	// Чтение сертификата CA
	clientCA, err := os.ReadFile("../certs/candy.tld/cert.pem")
	if err != nil {
		log.Fatalf("Reading cert failed: %v", err) // Завершаем программу в случае ошибки чтения сертификата
	}

	// Создаем пул сертификатов CA
	clientCAPool := x509.NewCertPool()
	clientCAPool.AppendCertsFromPEM(clientCA) // Добавляем сертификаты CA в пул

	// Настраиваем сервер
	server := http.Server{
		Addr: ":3333",
		TLSConfig: &tls.Config{
			ClientCAs:  clientCAPool,                   // Устанавливаем CA сертификаты для проверки клиента
			ClientAuth: tls.RequireAndVerifyClientCert, // Требуем и проверяем сертификат клиента
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				// Загружаем сертификат и ключ сервера
				c, err := tls.LoadX509KeyPair("../certs/candy.tld/cert.pem", "../certs/candy.tld/key.pem")
				if err != nil {
					log.Printf("Error loading key pair: %v", err) // Логируем ошибку загрузки ключа
					return nil, err
				}
				return &c, nil
			},
		},
	}

	log.Printf("Server starting on port 3333") // Логируем начало работы сервера
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("Server failed: %v", err) // Завершаем программу, если сервер не удалось запустить
	}
}
