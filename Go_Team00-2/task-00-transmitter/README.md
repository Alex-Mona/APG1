# Структура каталога проекта:
```
APG1-Bootcamp/Go_Team00-2/src/task-00-transmitter/
├── cmd/
│   └── server/
│       └── server.go
└── pkg/
    └── frequency/
        ├── frequency.pb.go
        └── frequency_grpc.pb.go
```


## bash: Linux
```
sudo apt update
sudo apt install -y protobuf-compiler
```

## bash: MacOS
```
brew install protobuf
```
### Пояснение по файлу .proto:
package frequency; — это имя пакета, используемое внутри .proto файла.
option go_package = "APG1-Bootcamp/Go_Team00-2/src/task-00-transmitter/"; — указывает на пакет Go, в котором будет использоваться сгенерированный код. Укажите здесь правильный путь для вашего проекта (обычно это путь внутри вашего модуля Go).

Пример кода frequency.proto:
```
syntax = "proto3";

package frequency;

option go_package = "APG1-Bootcamp/Go_Team00-2/src/task-00-transmitter/";

message FrequencyRequest {
    string session_id = 1;
    double frequency = 2;
    int64 timestamp = 3;
}

message FrequencyResponse {
    string session_id = 1;
    double frequency = 2;
    int64 timestamp = 3;
}

service FrequencyService {
    rpc GetFrequency (FrequencyRequest) returns (stream FrequencyResponse);
}
```
### Установите protoc-gen-go и protoc-gen-go-grpc из рекомендуемых пакетов:

## bash
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
Убедитесь, что пути к бинарным файлам protoc-gen-go и protoc-gen-go-grpc добавлены в ваш PATH. Обычно они находятся в $GOPATH/bin. Добавьте следующую строку в .bashrc или .zshrc, если это еще не сделано:

## bash
Например 
```vim .zshrc```
```
export PATH="$PATH:$(go env GOPATH)/bin"
```
### Теперь можно использовать protoc для генерации кода из .proto файла. Для этого выполните:

## bash Генерируем код:
```
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       frequency.proto
```
--go_out=. и --go-grpc_out=. указывают на директорию, куда будут сгенерированы файлы.
--go_opt=paths=source_relative и --go-grpc_opt=paths=source_relative обеспечивают генерацию файлов с относительными путями, что упрощает организацию файлов в проекте.
После этого у вас должны появиться два файла: frequency.pb.go и frequency_grpc.pb.go. Первый содержит код для работы с протоколом буфера, второй — код для работы с gRPC. 

**Эти файлы  frequency.pb.go и frequency_grpc.pb.go нужно перенести в дирректорию pkg/frequency, Иначе возникнет Ошибка could not import.


---

# Объяснение работы фукнций gRPC сервера.

### Импорты

```go
package main

import (
    "log"
    "math/rand"
    "net"
    "time"

    "github.com/google/uuid"
    "google.golang.org/grpc"
    pb "APG1-Bootcamp/Go_Team00-2/src/task-00-transmitter/pkg/frequency"
)
```

- **`"log"`**: Пакет для записи логов.
- **`"math/rand"`**: Пакет для генерации случайных чисел.
- **`"net"`**: Пакет для работы с сетевыми соединениями.
- **`"time"`**: Пакет для работы со временем.
- **`"github.com/google/uuid"`**: Пакет для генерации уникальных идентификаторов UUID.
- **`"google.golang.org/grpc"`**: Пакет для работы с gRPC, библиотека для реализации gRPC серверов и клиентов.
- **`pb "APG1-Bootcamp/Go_Team00-2/src/task-00-transmitter/pkg/frequency"`**: Импорт сгенерированного кода gRPC из пакета `frequency`. Это определение ваших сообщений и сервисов, основанных на вашем `.proto` файле.

### Константы

```go
const (
    port = ":50051"
)
```

- **`port`**: Константа, задающая порт, на котором будет слушать gRPC сервер. В данном случае это порт `50051`.

### Определение сервера

```go
type server struct {
    pb.UnimplementedFrequencyServiceServer
}
```

- **`server`**: Это структура, реализующая интерфейс gRPC сервера, который вы сгенерировали из `.proto` файла. В Go, структура реализует интерфейс, предоставляемый gRPC, добавляя конкретные методы для обработки вызовов.

### Реализация метода сервиса

```go
func (s *server) StreamFrequencies(_ *pb.Empty, stream pb.FrequencyService_StreamFrequenciesServer) error {
    sessionID := uuid.New().String()
    mean := rand.Float64()*20 - 10          // Выбор среднего значения из интервала [-10, 10]
    stdDev := rand.Float64()*1.2 + 0.3      // Выбор стандартного отклонения из интервала [0.3, 1.5]

    log.Printf("New session: %s, Mean: %f, StdDev: %f\n", sessionID, mean, stdDev)

    for {
        frequency := rand.NormFloat64()*stdDev + mean
        timestamp := time.Now().UTC().Unix()

        msg := &pb.FrequencyMessage{
            SessionId: sessionID,
            Frequency: frequency,
            Timestamp: timestamp,
        }

        if err := stream.Send(msg); err != nil {
            return err
        }

        time.Sleep(100 * time.Millisecond)
    }
}
```

- **`StreamFrequencies`**: Это метод, реализующий потоковый вызов сервиса gRPC. В вашем `.proto` файле должен быть определен метод `StreamFrequencies`, который принимает пустой запрос и возвращает поток сообщений.

  - **`sessionID := uuid.New().String()`**: Генерирует новый уникальный идентификатор с использованием библиотеки `uuid`.
  
  - **`mean := rand.Float64()*20 - 10`**: Вычисляет среднее значение для нормального распределения, которое находится в диапазоне от -10 до 10.
  
  - **`stdDev := rand.Float64()*1.2 + 0.3`**: Вычисляет стандартное отклонение для нормального распределения, которое находится в диапазоне от 0.3 до 1.5.
  
  - **`log.Printf(...)`**: Записывает информацию о новой сессии и ее параметрах в лог.

  - **Цикл `for`**: В бесконечном цикле генерируется случайное значение частоты из нормального распределения с заданными `mean` и `stdDev`, а также текущее время в формате Unix timestamp.
  
  - **`msg := &pb.FrequencyMessage{...}`**: Создает новое сообщение сгенерированных данных.
  
  - **`if err := stream.Send(msg); err != nil`**: Отправляет сообщение через поток. Если возникает ошибка при отправке, метод возвращает эту ошибку и завершает выполнение.
  
  - **`time.Sleep(100 * time.Millisecond)`**: Пауза между отправкой сообщений в 100 миллисекунд.

### Основная функция

```go
func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterFrequencyServiceServer(s, &server{})

    log.Printf("Server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
```

- **`net.Listen("tcp", port)`**: Создает TCP-листенер, который будет слушать на заданном порту.

- **`grpc.NewServer()`**: Создает новый gRPC сервер.

- **`pb.RegisterFrequencyServiceServer(s, &server{})`**: Регистрирует реализованный сервер (структура `server`) для обслуживания gRPC сервисов. Этот метод связывает ваш сервер с gRPC сервером.

- **`log.Printf("Server listening at %v", lis.Addr())`**: Логирует информацию о том, что сервер начал слушать на определенном адресе.

- **`s.Serve(lis)`**: Запускает сервер и начинает принимать входящие соединения. Если возникает ошибка, сервер завершает выполнение с логированием ошибки.

Этот код реализует gRPC сервер, который генерирует и отправляет случайные данные частот через потоковый сервис. Он включает генерацию случайных данных, настройку gRPC сервера и логирование ошибок. 

# Компилируем
```
go build -o transmitter server.go
./transmitter
```

---