package main

import (
	"IM-System/internal/client"
	"flag"
	"fmt"
)

func main() {
	// Command line parameter parsing
	flag.Parse()
	serverIp, serverPort := client.GetServerConfig()

	client := client.NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("\033[31m>>>>> Fail to connect server...\033[0m")
		return
	}

	fmt.Println("\033[32m>>>>> Connect server successfully\033[0m")

	// Block main goroutine to keep the client running
	fmt.Scanln()
}
