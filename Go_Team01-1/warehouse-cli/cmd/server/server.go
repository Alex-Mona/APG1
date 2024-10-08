package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Node struct {
	Address       string    `json:"address"`
	Port          string    `json:"port"`
	IsLeader      bool      `json:"is_leader"`
	LastHeartbeat time.Time // Время последнего сердцебиения
}

var (
	dataStore         = make(map[string]string) // Хранилище данных
	nodes             = make(map[string]*Node)  // Хранение информации об узлах
	replicationFactor = 2                       // Фактор репликации
	serverPort        string                    // Порт сервера
	mu                sync.Mutex                // Синхронизация доступа
	isLeader          = false                   // Флаг лидера
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide port")
	}
	serverPort = os.Args[1]
	fmt.Printf("Server started on port %s\n", serverPort)

	// Инициализация первого узла как лидера
	if len(os.Args) > 2 {
		// Подключение к другому узлу
		connectToCluster(os.Args[2])
	} else {
		// Являемся лидером
		isLeader = true
		nodes[serverPort] = &Node{
			Address:       "127.0.0.1",
			Port:          serverPort,
			IsLeader:      true,
			LastHeartbeat: time.Now(),
		}
	}

	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/heartbeat", handleHeartbeat)
	http.HandleFunc("/join", handleJoinCluster)

	// Старт Heartbeat для мониторинга других узлов
	go monitorNodes()

	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}

// Получение данных по ключу
func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(w, "Error: key cannot be empty", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	value, exists := dataStore[key]
	if !exists {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, value)
}

// Установка данных по ключу
func handleSet(w http.ResponseWriter, r *http.Request) {
	if !isLeader {
		http.Error(w, "Current node is not the leader", http.StatusForbidden)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Error: key cannot be empty", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	value := string(body)

	if _, err := uuid.Parse(key); err != nil {
		http.Error(w, "Error: Key is not a proper UUID4", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	dataStore[key] = value
	fmt.Printf("Set key %s with value %s\n", key, value)

	// Репликация данных на другие узлы
	replicateToNodes(key, value)

	fmt.Fprintf(w, "Created/Updated (2 replicas)")
}

// Удаление данных по ключу
func handleDelete(w http.ResponseWriter, r *http.Request) {
	if !isLeader {
		http.Error(w, "Current node is not the leader", http.StatusForbidden)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Error: key cannot be empty", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	delete(dataStore, key)
	fmt.Printf("Deleted key %s\n", key)

	// Репликация удаления на другие узлы
	replicateToNodes(key, "")

	fmt.Fprintf(w, "Deleted (2 replicas)")
}

// Обработка heartbeats
func handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	heartbeatResponse := struct {
		Leader        bool `json:"leader"`
		ReplicaFactor int  `json:"replica_factor"`
	}{
		Leader:        isLeader,
		ReplicaFactor: replicationFactor,
	}
	json.NewEncoder(w).Encode(heartbeatResponse)
}

// Присоединение к кластеру
func handleJoinCluster(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	port := r.URL.Query().Get("port")

	if address == "" || port == "" {
		http.Error(w, "Error: address and port cannot be empty", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	nodes[port] = &Node{
		Address:       address,
		Port:          port,
		IsLeader:      false,
		LastHeartbeat: time.Now(),
	}
	fmt.Printf("Node %s:%s joined the cluster\n", address, port)
}

// Подключение к кластеру
func connectToCluster(existingNodePort string) {
	fmt.Printf("Connecting to existing node at %s\n", existingNodePort)
	url := fmt.Sprintf("http://127.0.0.1:%s/join?address=127.0.0.1&port=%s", existingNodePort, serverPort)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error connecting to cluster: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("Joined the cluster successfully")
}

// Репликация данных на другие узлы с учетом фактора репликации
// Репликация данных на другие узлы с учетом фактора репликации
func replicateToNodes(key, value string) {
	replicatedCount := 0

	// Перебираем узлы для репликации
	for port := range nodes {
		if port == serverPort || replicatedCount >= replicationFactor {
			continue // Пропускаем текущий узел и завершаем, если достигнут фактор репликации
		}

		url := fmt.Sprintf("http://127.0.0.1:%s/set?key=%s", port, key)
		log.Printf("Replicating to node %s with URL %s", port, url) // Лог для отладки

		req, err := http.NewRequest("POST", url, strings.NewReader(value))
		if err != nil {
			log.Printf("Error creating request for node %s: %v", port, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Failed to replicate to node %s: %v, Status code: %d", port, err, resp.StatusCode)
		} else {
			replicatedCount++
			log.Printf("Successfully replicated to node %s", port)
		}
		defer resp.Body.Close()
	}

	if replicatedCount < replicationFactor {
		log.Printf("Warning: Only replicated to %d nodes (required: %d)", replicatedCount, replicationFactor)
	}
}

// Мониторинг состояния узлов
func monitorNodes() {
	for {
		time.Sleep(5 * time.Second)

		mu.Lock()
		leaderAlive := false
		for port, node := range nodes {
			if time.Since(node.LastHeartbeat) > 10*time.Second {
				fmt.Printf("Node %s is down\n", port)
				delete(nodes, port)

				// Если узел был лидером, его падение требует выбор нового лидера
				if node.IsLeader {
					fmt.Printf("Leader node %s has failed, electing new leader...\n", port)
					electNewLeader()
				}
			} else if node.IsLeader {
				leaderAlive = true
			}
		}
		// Если нет активного лидера, выбираем нового
		if !leaderAlive && isLeader {
			fmt.Println("Current node is now the leader.")
			isLeader = true
		}
		mu.Unlock()
	}
}

// Выбор нового лидера среди доступных узлов
func electNewLeader() {
	fmt.Println("Electing new leader...")
	for port, node := range nodes {
		if node != nil && port != serverPort {
			node.IsLeader = true
			nodes[port] = node
			isLeader = false
			fmt.Printf("Node %s has been elected as the new leader\n", port)
			return // Завершаем после выбора нового лидера
		}
	}
	fmt.Println("No available nodes to elect as leader.")
}
