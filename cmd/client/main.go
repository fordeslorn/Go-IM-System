package main

import (
	"IM-System/internal/client"
	"fmt"
)

func main() {
	client := client.NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("\033[31m>>>>> Fail to connect server...\033[0m")
		return
	}

	fmt.Println("\033[32m>>>>> Connect server successfully\033[0m")

	// Block main goroutine to keep the client running
	fmt.Scanln()
}
