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
		return 0, errors.New("wrong candy type") // Возвращаем ошибку, если тип конфет неизвестен
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
		w.WriteHeader(http.StatusBadRequest)
		response = struct {
			Error string `json:"error"`
		}{"Wrong candy type!"} // Возвращаем ошибку 400, если тип конфет неправильный
	} else if data.CandyCount < 0 {
		w.WriteHeader(http.StatusBadRequest)
		response = struct {
			Error string `json:"error"`
		}{"Negative candy count!"} // Возвращаем ошибку 400, если количество конфет отрицательное
	} else if candyPrice*data.CandyCount > data.Money {
		amount := candyPrice*data.CandyCount - data.Money
		w.WriteHeader(http.StatusPaymentRequired)
		response = struct {
			Error string `json:"error"`
		}{"You need " + strconv.Itoa(amount) + " more money!"} // Возвращаем ошибку 402, если недостаточно денег
	} else {
		change := data.Money - candyPrice*data.CandyCount
		w.WriteHeader(http.StatusOK)
		response = struct {
			Change int    `json:"change"`
			Thanks string `json:"thanks"`
		}{change, "Thank you!"} // Возвращаем оставшиеся деньги и благодарность
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

	// Читаем сертификат CA
	clientCA, err := os.ReadFile("../certs/candy.tld/cert.pem")
	if err != nil {
		log.Fatalf("Reading cert failed: %v", err) // Завершаем программу в случае ошибки чтения сертификата
	}

	// Создаем пул сертификатов CA
	clientCAPool := x509.NewCertPool()
	clientCAPool.AppendCertsFromPEM(clientCA)

	// Настраиваем сервер
	server := http.Server{
		Addr: ":3333",
		TLSConfig: &tls.Config{
			ClientCAs:  clientCAPool,                   // Устанавливаем сертификаты CA для проверки клиента
			ClientAuth: tls.RequireAndVerifyClientCert, // Требуем проверку клиентского сертификата
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				// Загружаем сертификат и ключ для сервера
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
