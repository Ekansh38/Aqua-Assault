package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	CONNECT := "localhost:1234"
	conn, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	go readFromServer(conn)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Client: ")
		message, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}

func readFromServer(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		fmt.Println(message)
	}
}
