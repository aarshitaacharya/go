package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

// ============================================================================
// 1. ARCHITECTURE STRUCTS
// ============================================================================

// Day 1: The Shared Memory Model (Mutex)
type MutexState struct {
	mu sync.RWMutex
	db map[string]string
}

// Day 2: The Communicating Sequential Processes Model (Channels)
type ChanState struct {
	db      map[string]string
	actions chan func()
}

// ============================================================================
// 2. MAIN APPLICATION ENTRYPOINT
// ============================================================================

func main() {
	// Initialize the Channel State for the running TCP server
	state := &ChanState{
		db:      make(map[string]string),
		actions: make(chan func()),
	}

	// Spawn the background manager to process incoming functions
	go runBackendManager(state)

	fmt.Println("Server running using Channel backend...")

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

func handleConnection(conn net.Conn, state *ChanState) {
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
		response := dispatchCommandChan(command, state)
		conn.Write([]byte(response))
	}
}

// ============================================================================
// 3. CHANNEL BACKEND DISPATCHER (Day 2 Approach)
// ============================================================================

func dispatchCommandChan(args []string, state *ChanState) string {
	if len(args) == 0 {
		return ""
	}
	switch args[0] {
	case "SET":
		if len(args) != 3 {
			return "ERR: Wrong number of arguments\n"
		}
		resChan := make(chan string)
		state.actions <- func() {
			state.db[args[1]] = args[2]
			resChan <- "OK\n"
		}
		return <-resChan

	case "GET":
		if len(args) != 2 {
			return "ERR: Wrong number of arguments\n"
		}
		resChan := make(chan string)
		state.actions <- func() {
			val, exists := state.db[args[1]]
			if !exists {
				resChan <- "(nil)\n"
			} else {
				resChan <- val + "\n"
			}
		}
		return <-resChan

	case "DEL":
		if len(args) != 2 {
			return "ERR: Wrong number of arguments\n"
		}
		resChan := make(chan string)
		state.actions <- func() {
			_, exists := state.db[args[1]]
			if !exists {
				resChan <- "0\n"
			} else {
				delete(state.db, args[1])
				resChan <- "1\n"
			}
		}
		return <-resChan

	case "EXISTS":
		if len(args) != 2 {
			return "ERR: Wrong number of arguments\n"
		}
		resChan := make(chan string)
		state.actions <- func() {
			_, exists := state.db[args[1]]
			if !exists {
				resChan <- "0\n"
			} else {
				resChan <- "1\n"
			}
		}
		return <-resChan

	case "CRASH":
		for i := 0; i < 100; i++ {
			go func(id int) {
				state.actions <- func() {
					state.db["collision_key"] = fmt.Sprintf("value-%d", id)
				}
			}(i)
		}
		return "Chaos unleashed\n"

	default:
		return "ERR: Command does not exist\n"
	}
}

func runBackendManager(state *ChanState) {
	for action := range state.actions {
		action()
	}
}

// ============================================================================
// 4. MUTEX BACKEND DISPATCHER (Day 1 Approach)
// ============================================================================

func dispatchCommandMutex(args []string, state *MutexState) string {
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
		state.mu.Lock()
		_, exists := state.db[args[1]]
		if !exists {
			state.mu.Unlock()
			return "0\n"
		}
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
		return "ERR: Command does not exist\n"
	}
}
