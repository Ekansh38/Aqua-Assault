package main

import (
	"bufio"
	"fmt"
	"net"
	// "os"
	"sync"
)

func removeClient(s []Client, i int) []Client {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func findClientBasedOnConn(s []Client, conn net.Conn) int {
	for i, client := range s {
		if client.conn == conn {
			return i
		}
	}
	return -1
}

type Client struct {
	conn net.Conn
	id   int
	x    int
	y    int
}

var (
	clients        []Client
	mu             sync.Mutex
	globalMessages []string
)

func messageToClient(conn net.Conn, message string) {
	conn.Write([]byte("Server: " + message))
}

func broadcastGlobalMessages() {
	for {
		mu.Lock()
		if len(globalMessages) > 0 {
			message := globalMessages[0]
			globalMessages = globalMessages[1:] // Remove the first message
			fmt.Println("Broadcasting:", message)

			// Send the message to all connected clients
			for _, client := range clients {
				_, err := client.conn.Write([]byte(message + "\n"))
				if err != nil {
					fmt.Printf("Error writing to client %d: %v\n", client.id, err)
				}
			}
		}
		mu.Unlock()
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		mu.Lock()
		i := findClientBasedOnConn(clients, conn)
		clients = removeClient(clients, i)
		mu.Unlock()
		conn.Close()
		fmt.Println("Client disconnected")
	}()

	mu.Lock()
	client := clients[findClientBasedOnConn(clients, conn)]

	globalMessages = append(globalMessages, fmt.Sprintf("Server: NCC%d POS%d,%d", client.id, client.x, client.y)) // NCC = New Client Connected

	mu.Unlock()

	go readFromClient(conn)

	for {
	}
}

func readFromClient(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}
		fmt.Println("Client:", message)
	}
}

func main() {
	PORT := ":1234"
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is running on port", PORT)

	go broadcastGlobalMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		mu.Lock()
		newClient := Client{conn: conn, id: len(clients), x: 0, y: 0}
		clients = append(clients, newClient)
		mu.Unlock()
		fmt.Println("New client connected")
		go handleConnection(conn)
	}
}
