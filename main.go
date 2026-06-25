// Today's task:
// Modify it to accept Set and Get commands using simple map[string]string

package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	db := make(map[string]string)

	fmt.Println("hello")

	listener, err := net.Listen("tcp", "localhost:6379")

	if err != nil {
		fmt.Println("Error setting up connection")
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connections")
			continue
		}

		go handleConnection(conn, db)
	}

}

func handleConnection(conn net.Conn, db map[string]string) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Client disconnected gracefully")
			return
		}

		input := string(buf[:n])
		command := strings.Fields(input)

		if command[0] == "SET" {
			db[command[1]] = command[2]
			_, err1 := conn.Write([]byte("OK\n"))
			if err1 != nil {
				fmt.Println(err1)
			}

		} else if command[0] == "GET" {
			val, exists := db[command[1]]
			if !exists {
				conn.Write([]byte("(nil)\n"))
			} else {
				conn.Write([]byte(val + "\n"))
			}
		} else {
			fmt.Println("Error, invalid value")
		}
	}
}
