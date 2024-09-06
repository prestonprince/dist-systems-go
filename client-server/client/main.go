package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("err connecting to server: ", err.Error())
		return
	}
	defer conn.Close()

	fmt.Println("connected to server. Type messages (or quit)")

	go readResponses(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "quit" {
			break
		}

		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("error sending msg: ", err.Error())
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading input: ", err.Error())
	}
}

func readResponses(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println("server response: ", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("error reading from server: ", err.Error())
	}
}
