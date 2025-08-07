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

	// start a goroutine to deal with server responses
	go client.DealResponse()

	fmt.Println("\033[32m>>>>> Connect server successfully\033[0m")

	client.Run()
}
