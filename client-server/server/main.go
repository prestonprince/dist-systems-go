package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Starting server...")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}
		defer conn.Close()

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("New Client Connected: ", conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("received message: %s\n", message)

		conn.Write([]byte(strings.ToUpper(message) + "\n"))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from connection: ")
	}

	fmt.Println("Client disconnected: ", conn.RemoteAddr().String())
}
