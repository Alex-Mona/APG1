package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"

)

func crawlWeb(ctx context.Context, urls chan string) chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8) // Ограничиваем количество параллельных запросов до 8

	for url := range urls {
		select {
		case <-ctx.Done():
			fmt.Println("Crawling stopped.")
			close(out)
			return out
		default:
			wg.Add(1)
			sem <- struct{}{} // Добавляем горутину в семафор
			go func(url string) {
				defer wg.Done()
				defer func() { <-sem }() // Освобождаем семафор

				resp, err := http.Get(url)
				if err != nil {
					fmt.Println("Error fetching URL:", err)
					return
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error reading response:", err)
					return
				}
				out <- string(body)
			}(url)
		}
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	urls := make(chan string, 5)
	ctx, cancel := context.WithCancel(context.Background())

	// Запускаем обработчик для Ctrl+C
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()

	urlList := []string{
		"https://example.com",
		"https://golang.org",
	}

	// Добавляем URL'ы в канал
	go func() {
		for _, url := range urlList {
			urls <- url
		}
		close(urls)
	}()

	// Читаем результаты из канала
	for body := range crawlWeb(ctx, urls) {
		fmt.Println("Page body:", body[:60]) // Выводим первые 60 символов тела страницы
	}
}
