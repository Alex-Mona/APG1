package main

import (
	"log" // Пакет для записи логов
	"math/rand" // Пакет для генерации случайных чисел
	"net" // Пакет для работы с сетевыми соединениями
	"time" // Пакет для работы со временем

	pb "APG1-Bootcamp/Go_Team00-2/src/task-01-anomaly-detection/pkg/frequency" // Импорт сгенерированного gRPC-кода

	"github.com/google/uuid" // Пакет для генерации уникальных идентификаторов UUID
	"google.golang.org/grpc" // Пакет для работы с gRPC
)

const (
    port = ":50051" // Порт, на котором будет работать сервер
)

// Определение структуры сервера, реализующего интерфейс gRPC
type server struct {
    pb.UnimplementedFrequencyServiceServer // Встраивание не реализованного сервиса для частичного использования
}

// Реализация метода StreamFrequencies, который будет использоваться для обработки потоковых запросов
func (s *server) StreamFrequencies(_ *pb.Empty, stream pb.FrequencyService_StreamFrequenciesServer) error {
    // Генерация нового уникального идентификатора для сессии
    sessionID := uuid.New().String()
    
    // Генерация случайного среднего значения из интервала [-10, 10]
    mean := rand.Float64()*20 - 10
    
    // Генерация случайного стандартного отклонения из интервала [0.3, 1.5]
    stdDev := rand.Float64()*1.2 + 0.3

    // Логирование информации о новой сессии и ее параметрах
    log.Printf("New session: %s, Mean: %f, StdDev: %f\n", sessionID, mean, stdDev)

    // Бесконечный цикл для отправки частот в потоковом режиме
    for {
        // Генерация случайного значения частоты на основе нормального распределения
        frequency := rand.NormFloat64()*stdDev + mean
        
        // Получение текущего времени в формате Unix timestamp
        timestamp := time.Now().UTC().Unix()

        // Создание нового сообщения с частотой и метаданными
        msg := &pb.FrequencyMessage{
            SessionId: sessionID,
            Frequency: frequency,
            Timestamp: timestamp,
        }

        // Отправка сообщения через поток. Если возникает ошибка, метод завершает работу
        if err := stream.Send(msg); err != nil {
            return err
        }

        // Пауза между отправками сообщений в 100 миллисекунд
        time.Sleep(100 * time.Millisecond)
    }
}

// Основная функция, запускающая сервер
func main() {
    // Создание TCP-листенера для прослушивания на указанном порту
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err) // Логирование ошибки при создании листенера
    }

    // Создание нового gRPC сервера
    s := grpc.NewServer()
    
    // Регистрация сервера для обработки gRPC вызовов
    pb.RegisterFrequencyServiceServer(s, &server{})

    // Логирование информации о том, что сервер начал прослушивание на указанном адресе
    log.Printf("Server listening at %v", lis.Addr())
    
    // Запуск сервера для обработки входящих соединений. Если возникает ошибка, сервер завершает работу
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err) // Логирование ошибки при запуске сервера
    }
}
