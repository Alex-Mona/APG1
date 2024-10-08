package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	currentNode = "127.0.0.1:8765"                                               // Узел по умолчанию
	nodes       = []string{"127.0.0.1:8765", "127.0.0.1:8766", "127.0.0.1:8767"} // Список узлов для переключения
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Connected to the database")
	fmt.Println("Available commands: GET <key>, SET <key> <value>, DELETE <key>, JOIN <address> <port>")

	for {
		fmt.Print("> ")
		scanner.Scan()
		command := scanner.Text()

		if strings.HasPrefix(command, "GET") {
			key := strings.TrimSpace(strings.TrimPrefix(command, "GET"))
			if key == "" {
				fmt.Println("Error: key cannot be empty")
				continue
			}
			get(key)
		} else if strings.HasPrefix(command, "SET") {
			parts := strings.SplitN(strings.TrimPrefix(command, "SET "), " ", 2)
			if len(parts) != 2 {
				fmt.Println("Usage: SET <key> <value>")
				continue
			}
			key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			if key == "" || value == "" {
				fmt.Println("Error: key and value cannot be empty")
				continue
			}
			set(key, value)
		} else if strings.HasPrefix(command, "DELETE") {
			key := strings.TrimSpace(strings.TrimPrefix(command, "DELETE"))
			if key == "" {
				fmt.Println("Error: key cannot be empty")
				continue
			}
			delete(key)
		} else if strings.HasPrefix(command, "JOIN") {
			parts := strings.Split(strings.TrimPrefix(command, "JOIN "), " ")
			if len(parts) != 2 {
				fmt.Println("Usage: JOIN <address> <port>")
				continue
			}
			address, port := parts[0], parts[1]
			joinCluster(address, port)
		} else {
			fmt.Println("Unknown command")
		}
	}
}

// GET запрос
func get(key string) {
	resp, err := http.Get(fmt.Sprintf("http://%s/get?key=%s", currentNode, key))
	if err != nil {
		fmt.Println("Error:", err)
		switchNode() // Переключение на другой узел при ошибке
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Println(string(body))
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Error: %s\n", string(body))
	}
}

// SET запрос
func set(key, value string) {
	escapedKey := url.QueryEscape(key)
	url := fmt.Sprintf("http://%s/set?key=%s", currentNode, escapedKey)

	req, err := http.NewRequest("POST", url, strings.NewReader(value))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		switchNode() // Переключение на другой узел при ошибке
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Created/Updated (2 replicas)")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Error setting value: %s\n", string(body))
	}
}

// DELETE запрос
func delete(key string) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/delete?key=%s", currentNode, key), nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		switchNode() // Переключение на другой узел при ошибке
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Deleted (2 replicas)")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Error: %s\n", string(body))
	}
}

// JOIN команда для присоединения к кластеру
func joinCluster(address, port string) {
	url := fmt.Sprintf("http://%s/join?address=%s&port=%s", currentNode, address, port)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error joining cluster:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Node joined the cluster successfully")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Error joining cluster: %s\n", string(body))
	}
}

// Переключение на следующий узел в случае ошибки
func switchNode() {
	for _, node := range nodes {
		if node != currentNode {
			currentNode = node
			fmt.Printf("Switched to node: %s\n", currentNode)
			return
		}
	}
	fmt.Println("No available nodes")
}
