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
		response := dispatchCommand(command, db)
		conn.Write(([]byte(response)))
	}
}

func dispatchCommand(args []string, db map[string]string) string {
	if len(args) == 0 {
		return ""
	}
	switch args[0] {
	case "SET":
		if len(args) != 3 {
			return "ERR: Wrong number of arguments\n"
		}
		db[args[1]] = args[2]
		return "OK\n"

	case "GET":
		if len(args) != 2 {
			return "ERR: Wrong number of arguments\n"
		}
		val, exists := db[args[1]]
		if !exists {
			return "(nil)\n"
		}
		return val + "\n"

	case "DEL":
		if len(args) != 2 {
			return "ERR: Wrong number of arguments\n"
		}
		_, exists := db[args[1]]
		if !exists {
			return "0\n"
		}
		delete(db, args[1])
		return "1\n"

	case "EXISTS":
		if len(args) != 2 {
			return "ERR: Wrong number of arguments\n"
		}
		_, exists := db[args[1]]
		if !exists {
			return "0\n"
		}
		return "1\n"

	case "CRASH":
		for i := 0; i < 100; i++ {
			go func(id int) {
				db["collision_key"] = fmt.Sprintf("value-%d", id)
			}(i)
		}
		return "Chaos unleashed\n"

	default:
		return "ERR: Command does not exists"
	}
}
