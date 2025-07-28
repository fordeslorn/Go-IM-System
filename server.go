package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// Create a new server instance
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (s *Server) Handler(conn net.Conn) {
	// Handle the connection
	fmt.Println("\033[32mConnection established successfully!\033[0m")
	fmt.Printf("\033[34mClient connected from:\033[0m [%s]%s\n", conn.RemoteAddr().Network(), conn.RemoteAddr().String())
}

// Start the server
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	// close connections
	defer listener.Close()

	for {
		// accept connections
		conn, err := listener.Accept() //conn is the socket
		if err != nil {
			fmt.Println("listener.Accept error:", err)
			continue
		}
		// handle connections
		go s.Handler(conn)
	}
}
