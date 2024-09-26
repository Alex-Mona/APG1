package main

import (
	"context"
	"testing"
	"time"
)

// Заглушка для имитации загрузки страниц
func TestCrawlWeb(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	urls := make(chan string, 3)
	urls <- "https://example.com"
	urls <- "https://golang.org"
	close(urls)

	// Запуск функции crawlWeb
	results := crawlWeb(ctx, urls)

	// Ожидание результата
	count := 0
	for range results {
		count++
	}

	// Ожидаем, что хотя бы 2 страницы были загружены
	if count < 2 {
		t.Errorf("Expected at least 2 pages to be crawled, got %d", count)
	}
}

// Тест остановки по Ctrl+C (симуляция отмены)
func TestCrawlWebCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	urls := make(chan string, 5)

	// Добавляем несколько URL для загрузки
	go func() {
		urls <- "https://example.com"
		time.Sleep(500 * time.Millisecond)
		cancel() // Прерываем выполнение
		close(urls)
	}()

	// Проверяем, что процесс был прерван
	results := crawlWeb(ctx, urls)
	select {
	case <-results:
		// Ожидаем, что результат придет до того, как тест завершится
	case <-time.After(3 * time.Second):
		t.Errorf("Test timed out, cancellation may not be working correctly")
	}
}
