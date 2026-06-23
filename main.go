// Write a program that prints "hello".
// Then write a TCP server that listens on
// port 6379 and prints whatever bytes come in.
// Use netcat to test: nc localhost 6379.

// Step 1: Print hello
// Step 2: Open a TCP server on port 6379, and handle errors
// Step 3: Infinite loop to listens
// Step 4: Create a buffer to store bytes
// Step 5: Print those bytes

package main

import (
	"fmt"
	"net"
)

func main() {
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

		handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(buf[:n]))
	}
}
