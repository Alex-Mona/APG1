package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Data представляет структуру данных для запроса
type Data struct {
	Money      int    `json:"money"`      // Количество денег
	CandyType  string `json:"candyType"`  // Тип конфет
	CandyCount int    `json:"candyCount"` // Количество конфет
}

// isFlagPassed проверяет, был ли установлен флаг с указанным именем
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	// Определение флагов командной строки
	candyType := flag.String("k", "", "two-letter abbreviation for the candy type")
	candyCount := flag.Int("c", 0, "count of candy to buy")
	money := flag.Int("m", 0, "amount of money given to machine")
	flag.Parse()

	// Проверка, были ли переданы все необходимые флаги
	if !isFlagPassed("k") || !isFlagPassed("c") || !isFlagPassed("m") {
		log.Fatalln("Missing required arguments")
	}

	// Создание структуры данных для запроса
	data := Data{
		Money:      *money,
		CandyType:  *candyType,
		CandyCount: *candyCount,
	}

	// Сериализация данных в формат JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("Error encoding data:", err) // Завершаем программу, если ошибка сериализации данных
	}

	// Загрузка сертификата корневого удостоверяющего центра
	rootCA, err := os.ReadFile("../certs/candy.tld/cert.pem")
	if err != nil {
		log.Fatalf("Error reading root CA certificate: %v", err)
	}
	rootCAPool := x509.NewCertPool()
	rootCAPool.AppendCertsFromPEM(rootCA)

	// Создание клиента HTTP с TLS настройками
	client := http.Client{
		Timeout: 15 * time.Second, // Устанавливаем таймаут на 15 секунд
		Transport: &http.Transport{
			IdleConnTimeout: 10 * time.Second, // Таймаут для неактивных соединений
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAPool, // Устанавливаем корневой сертификат для проверки сервера
				GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
					// Загружаем сертификат и ключ клиента
					c, err := tls.LoadX509KeyPair("../certs/candy.tld/cert.pem", "../certs/candy.tld/key.pem")
					if err != nil {
						log.Printf("Error loading key pair: %v", err)
						return nil, err
					}
					return &c, nil
				},
			},
		},
	}

	// Создание и отправка HTTP запроса
	req, err := http.NewRequest(http.MethodPost, "https://candy.tld:3333/buy_candy", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln("Error creating request:", err) // Завершаем программу, если ошибка создания запроса
	}

	req.Header.Set("Content-Type", "application/json") // Устанавливаем заголовок типа контента

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error sending request:", err) // Завершаем программу, если ошибка отправки запроса
	}
	defer resp.Body.Close() // Закрываем тело ответа после обработки

	// Копирование тела ответа в стандартный вывод
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		log.Fatalln("Error reading response body:", err) // Завершаем программу, если ошибка чтения ответа
	}
}

// curl -s --key ../certs/candy.tld/key.pem --cert ../certs/candy.tld/cert.pem --cacert ../certs/minica.pem -XPOST -H "Content-Type: application/json" -d '{"candyType": "NT", "candyCount": 2, "money": 34}' "https://candy.tld:3333/buy_candy"
