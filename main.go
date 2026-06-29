package main

import (
	"fmt"
	"net"
	"strings"
)

type State struct {
	db      map[string]string
	actions chan func()
}

func main() {
	state := &State{
		db:      make(map[string]string),
		actions: make(chan func()),
	}
	go runBackendManager(state)

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
		return "Chaos"

	default:
		return "ERR: Command does not exists"
	}
}

func runBackendManager(state *State) {
	for action := range state.actions {
		action()
	}
}
