// Today's task:
// Modify it to accept Set and Get commands using simple map[string]string

package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type State struct {
	mu sync.RWMutex
	db map[string]string
}

func main() {
	state := &State{
		db: make(map[string]string),
	}

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

		go handleConnection(conn, state)
	}

}

func handleConnection(conn net.Conn, state *State) {
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
		response := dispatchCommand(command, state)
		conn.Write(([]byte(response)))
	}
}

func dispatchCommand(args []string, state *State) string {
	if len(args) == 0 {
		return ""
	}
	switch args[0] {
	case "SET":
		if len(args) != 3 {
			return "ERR: Wrong number of arguments\n"
		}
		state.mu.Lock()
		state.db[args[1]] = args[2]
		state.mu.Unlock()
		return "OK\n"

	case "GET":
		if len(args) != 2 {
			return "ERR: Wrong number of arguments\n"
		}
		state.mu.RLock()
		val, exists := state.db[args[1]]
		state.mu.RUnlock()
		if !exists {
			return "(nil)\n"
		}
		return val + "\n"

	case "DEL":
		if len(args) != 2 {
			return "ERR: Wrong number of arguments\n"
		}
		_, exists := state.db[args[1]]
		if !exists {
			return "0\n"
		}
		state.mu.Lock()
		delete(state.db, args[1])
		state.mu.Unlock()
		return "1\n"

	case "EXISTS":
		if len(args) != 2 {
			return "ERR: Wrong number of arguments\n"
		}
		state.mu.RLock()
		_, exists := state.db[args[1]]
		state.mu.RUnlock()
		if !exists {
			return "0\n"
		}
		return "1\n"

	case "CRASH":
		for i := 0; i < 100; i++ {
			go func(id int) {
				state.mu.Lock()
				state.db["collision_key"] = fmt.Sprintf("value-%d", id)
				state.mu.Unlock()
			}(i)
		}
		return "Chaos unleashed\n"

	default:
		return "ERR: Command does not exists"
	}
}
