package client

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// link to server
	conn, err := net.Dial("tcp", net.JoinHostPort(serverIp, fmt.Sprintf("%d", serverPort)))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn

	return client
}

var serverIp string
var serverPort int

func GetServerConfig() (string, int) {
	return serverIp, serverPort
}

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server IP address(default 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "set server Port(default 8888)")
}
