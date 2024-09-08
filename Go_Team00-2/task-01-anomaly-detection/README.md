# gRPC Client для Обнаружения Аномалий

Этот клиент на Go подключается к gRPC серверу и обрабатывает поток данных частот для обнаружения аномалий.

## Импорт необходимых пакетов

```go
import (
    "context"
    "flag"
    "log"
    "math"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "APG1-Bootcamp/Go_Team00-2/src/task-01-anomaly-detection/pkg/frequency"
)
```

- **`"context"`**: Пакет для управления временем жизни операций.
- **`"flag"`**: Пакет для обработки командных флагов.
- **`"log"`**: Пакет для логирования информации и ошибок.
- **`"math"`**: Пакет для математических функций, таких как вычисление квадратного корня.
- **`"google.golang.org/grpc"`**: Пакет для gRPC функциональности.
- **`"google.golang.org/grpc/credentials/insecure"`**: Пакет для создания небезопасного транспортного соединения (не рекомендуется для продакшн-окружения).

## Константы

```go
const (
    defaultAddress = "localhost:50051"
    defaultK       = 2.0
)
```

- **`defaultAddress`**: Адрес сервера по умолчанию.
- **`defaultK`**: Значение коэффициента обнаружения аномалий по умолчанию.

## Функция `main`

### Чтение параметров командной строки

```go
address := flag.String("address", defaultAddress, "Server address")
k := flag.Float64("k", defaultK, "Anomaly detection coefficient")
flag.Parse()
```

- **`address`** и **`k`**: Параметры задаются через командные флаги.
- **`flag.Parse()`**: Анализирует флаги командной строки.

### Создание контекста

```go
ctx := context.Background()
```

- **`ctx := context.Background()`**: Создает базовый контекст для gRPC операций.

### Создание gRPC клиента

```go
conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}
defer conn.Close()

client := pb.NewFrequencyServiceClient(conn)
```

- **`conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))`**: Устанавливает соединение с сервером, используя небезопасные транспортные креденшелы.
- **`defer conn.Close()`**: Закрывает соединение при завершении работы функции `main`.
- **`client := pb.NewFrequencyServiceClient(conn)`**: Создает новый gRPC клиент для общения с сервисом частот.

### Вызов метода `StreamFrequencies`

```go
stream, err := client.StreamFrequencies(ctx, &pb.Empty{})
if err != nil {
    log.Fatalf("Failed to start streaming: %v", err)
}
```

- **`stream, err := client.StreamFrequencies(ctx, &pb.Empty{})`**: Отправляет запрос на начало стриминга частот.
- **`if err != nil { log.Fatalf("Failed to start streaming: %v", err) }`**: Логирует ошибку, если стриминг не удалось начать.

### Обработка полученных данных

```go
for {
    msg, err := stream.Recv()
    if err != nil {
        log.Fatalf("Stream error: %v", err)
    }

    count++
    frequency := msg.GetFrequency()
    sum += frequency
    sumSq += frequency * frequency

    // Рассчитываем среднее и стандартное отклонение
    mean = sum / float64(count)
    variance := (sumSq / float64(count)) - (mean * mean)
    if variance < 0 {
        variance = 0
    }
    stdDev = math.Sqrt(variance)

    log.Printf("Received frequency: %f, Mean: %f, StdDev: %f", frequency, mean, stdDev)

    // Обнаружение аномалий
    if math.Abs(frequency-mean) > *k*stdDev {
        log.Printf("Anomaly detected! Frequency: %f, Mean: %f, StdDev: %f", frequency, mean, stdDev)
    }
}
```

- **В бесконечном цикле `for`**:
  - **`msg, err := stream.Recv()`**: Получает сообщение из потока.
  - **`count++`**: Увеличивает счетчик полученных сообщений.
  - **`frequency := msg.GetFrequency()`**: Извлекает частоту из сообщения.
  - **`sum += frequency` и `sumSq += frequency * frequency`**: Обновляет сумму и сумму квадратов частот для вычисления среднего и дисперсии.
  
### Расчет статистических параметров

```go
mean = sum / float64(count)
variance := (sumSq / float64(count)) - (mean * mean)
if variance < 0 {
    variance = 0
}
stdDev = math.Sqrt(variance)
```

- **`mean = sum / float64(count)`**: Вычисляет среднее значение.
- **`variance := (sumSq / float64(count)) - (mean * mean)`**: Вычисляет дисперсию.
- **`if variance < 0 { variance = 0 }`**: Гарантирует, что дисперсия не будет отрицательной.
- **`stdDev = math.Sqrt(variance)`**: Вычисляет стандартное отклонение.

### Логирование и обнаружение аномалий

```go
log.Printf("Received frequency: %f, Mean: %f, StdDev: %f", frequency, mean, stdDev)
if math.Abs(frequency-mean) > *k*stdDev {
    log.Printf("Anomaly detected! Frequency: %f, Mean: %f, StdDev: %f", frequency, mean, stdDev)
}
```

- **`log.Printf("Received frequency: %f, Mean: %f, StdDev: %f", frequency, mean, stdDev)`**: Логирует полученные данные.
- **`if math.Abs(frequency-mean) > *k*stdDev`**: Проверяет, является ли частота аномалией, если отклонение превышает заданное количество стандартных отклонений.
- **`log.Printf("Anomaly detected! Frequency: %f, Mean: %f, StdDev: %f", frequency, mean, stdDev)`**: Логирует информацию об обнаруженной аномалии.

# Компилируем
```
go build -o anomaly_detector client.go
./anomaly_detector -k=3.0
```
---

Этот README.md должен помочь вам понять, как работает клиент gRPC и как он взаимодействует с сервером для обнаружения аномалий в данных частот.