package main

import (
	"IM-System/internal/server"
	"fmt"
)

func main() {
	fmt.Println("\033[34mStarting the server...\033[0m")
	server := server.NewServer("localhost", 8888)
	server.Start()
}
