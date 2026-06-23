// Write a program that prints "hello".
// Then write a TCP server that listens on
// port 6379 and prints whatever bytes come in.
// Use netcat to test: nc localhost 6379.

// Step 1: Print hello
// Step 2: Create a TCP server that listens on port 6379 and handle errors
// Step 3: Infinite waiting loop to pause and wait
// Step 4: Set up a buffer to read incoming data
// Step 5: Print the received data to the console

package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("hello")
	listener, err := net.Listen("tcp", "localhost:6379")

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer listener.Close()
	fmt.Println("Server successfully listening on port 6379")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection", err)
			continue
		}

		fmt.Println("Client connected!")

		conn.Close()
	}
}
