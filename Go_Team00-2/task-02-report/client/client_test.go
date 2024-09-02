package main

import (
    "context"
    "log"
    "math"
    "net"
    "testing"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/test/bufconn"
    pb "APG1-Bootcamp/Go_Team00-2/src/task-02-report/pkg/frequency" // Путь к вашему сгенерированному пакету
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

const bufSize = 1024 * 1024 // Размер буфера для буферного соединения

var lis *bufconn.Listener

// Реализация сервера gRPC, который будет использоваться в тестах
type server struct {
    pb.UnimplementedFrequencyServiceServer
}

// Метод StreamFrequencies, который отправляет тестовые данные клиенту
func (s *server) StreamFrequencies(req *pb.Empty, stream pb.FrequencyService_StreamFrequenciesServer) error {
    // Генерация простых данных для стрима
    for i := 0; i < 10; i++ {
        response := &pb.FrequencyMessage{
            Frequency: float64(i) + 10.0, // Простая генерация частот
            SessionId: "test-session",
            Timestamp: time.Now().Unix(),
        }
        if err := stream.Send(response); err != nil {
            return err // Возвращаем ошибку, если не удалось отправить данные
        }
        time.Sleep(500 * time.Millisecond) // Задержка для имитации реального стрима
    }
    return nil
}

func init() {
    lis = bufconn.Listen(bufSize) // Создаем буферное соединение
    s := grpc.NewServer()         // Создаем новый gRPC сервер
    pb.RegisterFrequencyServiceServer(s, &server{}) // Регистрируем сервер gRPC
    go func() {
        if err := s.Serve(lis); err != nil {
            log.Fatalf("Server exited with error: %v", err) // Логируем ошибку, если сервер завершился с ошибкой
        }
    }()
}

// bufDialer создает подключение к буферному серверу
func bufDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial() // Возвращаем соединение с буферным сервером
}

// Тестовая функция для проверки основной логики
func TestMainFunction(t *testing.T) {
    // Эмулируем подключение к gRPC серверу
    conn, err := grpc.DialContext(
        context.Background(), // Контекст для соединения
        "bufnet",            // Название сети
        grpc.WithContextDialer(bufDialer), // Диалер для подключения
        grpc.WithTransportCredentials(insecure.NewCredentials()), // Не безопасное соединение
    )
    if err != nil {
        t.Fatalf("Failed to dial bufnet: %v", err) // Логируем ошибку, если не удалось установить соединение
    }
    defer conn.Close() // Закрываем соединение после завершения теста

    client := pb.NewFrequencyServiceClient(conn) // Создаем новый клиент gRPC

    // Подключаемся к тестовой базе данных SQLite (вместо PostgreSQL)
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err) // Логируем ошибку, если не удалось подключиться к базе данных
    }

    // Мигрируем схему базы данных
    if err := db.AutoMigrate(&Anomaly{}); err != nil {
        t.Fatalf("Failed to migrate database schema: %v", err) // Логируем ошибку, если не удалось выполнить миграцию
    }

    // Инициализируем контекст с таймаутом
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // Отменяем контекст после завершения теста

    // Стриминг частот
    stream, err := client.StreamFrequencies(ctx, &pb.Empty{})
    if err != nil {
        t.Fatalf("Failed to start streaming: %v", err) // Логируем ошибку, если не удалось начать стриминг
    }

    var sum float64
    var sumSq float64
    var count int
    var mean float64
    var stdDev float64
    k := 2.0 // Коэффициент для определения аномалий

    for i := 0; i < 10; i++ { // Получаем первые 10 сообщений
        msg, err := stream.Recv()
        if err != nil {
            t.Fatalf("Stream error: %v", err) // Логируем ошибку, если возникла проблема с чтением из потока
        }

        count++
        frequency := msg.GetFrequency() // Получаем частоту из сообщения
        sum += frequency // Обновляем сумму частот
        sumSq += frequency * frequency // Обновляем сумму квадратов частот

        // Рассчитываем среднее значение и стандартное отклонение
        mean = sum / float64(count)
        variance := (sumSq / float64(count)) - (mean * mean)
        if variance < 0 {
            variance = 0 // Гарантируем, что дисперсия не отрицательная
        }
        stdDev = math.Sqrt(variance) // Вычисляем стандартное отклонение

        t.Logf("Received frequency: %f, Mean: %f, StdDev: %f", frequency, mean, stdDev)

        // Проверка на аномалии
        if math.Abs(frequency-mean) > k*stdDev {
            t.Logf("Anomaly detected! Frequency: %f, Mean: %f, StdDev: %f", frequency, mean, stdDev)

            anomaly := Anomaly{
                Frequency: frequency,
                Mean:      mean,
                StdDev:    stdDev,
                DetectedAt: time.Now(),
            }
            if err := db.Create(&anomaly).Error; err != nil {
                t.Fatalf("Failed to record anomaly: %v", err) // Логируем ошибку, если не удалось записать аномалию в базу данных
            }
        }
    }

    // Проверяем количество записанных аномалий
    var anomaliesCount int64
    if err := db.Model(&Anomaly{}).Count(&anomaliesCount).Error; err != nil {
        t.Fatalf("Failed to count anomalies: %v", err) // Логируем ошибку, если не удалось подсчитать количество аномалий
    }

    t.Logf("Total anomalies recorded: %d", anomaliesCount) // Логируем общее количество записанных аномалий
}
