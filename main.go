package main

import "fmt"

func main() {
	fmt.Println("\033[34mStarting the server...\033[0m")
	server := NewServer("localhost", 8888)
	server.Start()
}
