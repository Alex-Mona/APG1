package main

import (
    "context"
    "testing"
    "net"
    "log"

    pb "APG1-Bootcamp/Go_Team00-2/src/task-02-report/pkg/frequency" // Импорт сгенерированного gRPC-кода
    "google.golang.org/grpc" // Пакет для работы с gRPC
    "google.golang.org/grpc/test/bufconn" // Пакет для тестирования gRPC с использованием буферизованных соединений
)

const bufSize = 1024 * 1024 // Размер буфера для буферизованного соединения

var lis *bufconn.Listener // Определяем переменную для буферизованного соединения

// Функция инициализации, которая выполняется перед запуском тестов
func init() {
    lis = bufconn.Listen(bufSize) // Создаем буферизованное соединение с заданным размером буфера
    s := grpc.NewServer() // Создаем новый gRPC сервер
    pb.RegisterFrequencyServiceServer(s, &server{}) // Регистрируем наш сервер в gRPC
    go func() {
        // Запускаем сервер для обработки соединений, если возникает ошибка - логируем и завершаем выполнение
        if err := s.Serve(lis); err != nil {
            log.Fatalf("Server exited with error: %v", err)
        }
    }()
}

// Функция, которая будет использоваться для создания соединения с сервером через буферизованное соединение
func bufDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial() // Возвращаем соединение, установленное через буферизованный слушатель
}

// Основной тест, который проверяет работу метода StreamFrequencies
func TestStreamFrequencies(t *testing.T) {
    ctx := context.Background() // Создаем контекст для gRPC вызовов

    // Устанавливаем соединение с сервером через буферизованное соединение
    conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to dial bufnet: %v", err) // Логируем и завершаем тест, если не удалось установить соединение
    }
    defer conn.Close() // Закрываем соединение по завершении теста

    client := pb.NewFrequencyServiceClient(conn) // Создаем клиента для общения с нашим gRPC сервером

    // Вызываем метод StreamFrequencies и получаем поток
    stream, err := client.StreamFrequencies(ctx, &pb.Empty{})
    if err != nil {
        t.Fatalf("Failed to call StreamFrequencies: %v", err) // Логируем и завершаем тест, если вызов метода завершился ошибкой
    }

    // Пробуем получить несколько сообщений из потока
    for i := 0; i < 10; i++ {
        msg, err := stream.Recv() // Получаем сообщение из потока
        if err != nil {
            t.Fatalf("Failed to receive a message: %v", err) // Логируем и завершаем тест, если не удалось получить сообщение
        }
        t.Logf("Received message: SessionID: %s, Frequency: %f, Timestamp: %d", msg.GetSessionId(), msg.GetFrequency(), msg.GetTimestamp()) // Логируем полученное сообщение

        // Проверяем, что частота не равна нулю
        if msg.GetFrequency() == 0 {
            t.Errorf("Expected a non-zero frequency") // Сообщаем об ошибке, если частота равна нулю
        }
    }
}

// Функция TestMain запускает тесты и может использоваться для выполнения предварительных настроек
func TestMain(m *testing.M) {
    // Используем функцию `init()` для запуска gRPC-сервера
    m.Run() // Запуск всех тестов
}
