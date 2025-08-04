package main

import (
	"fmt"
	"net"
	"sync"
	"io"
)

type Server struct {
	Ip   string
	Port int

	// Online users
	OnlineUserMap map[string]*User
	maplock       sync.RWMutex

	// Message channel for broadcasting messages
	MessageChan chan string
}

// Create a new server instance
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:            ip,
		Port:          port,
		OnlineUserMap: make(map[string]*User),
		MessageChan:   make(chan string, 100),
	}
	return server
}

func (s *Server) ListenMessager() {
	for {
		msg := <-s.MessageChan

		s.maplock.RLock()
		for _, user := range s.OnlineUserMap {
			user.C <- msg // Send the message to each user's channel
		}
		s.maplock.RUnlock()
	}
}

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + ":" + msg

	s.MessageChan <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	// Handle the connection
	fmt.Printf("\033[32mClient connection established successfully from:\033[0m [%s]%s\n", conn.RemoteAddr().Network(), conn.RemoteAddr().String())

	user := NewUser(conn, s)

	user.Online() 

	// Accept messages from the user
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				fmt.Printf("[%s]%s \033[34mhas disconnected\033[0m\n", conn.RemoteAddr().Network(), user.Name)
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn.Read error:", err)
				return
			}

			// get the message
			msg := string(buf[:n-1]) // Exclude the '\n' character

			// User handle message
			user.DoMessage(msg)
		}
	}()

	// Block the handler
	select {}
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

	// start listening for messages
	go s.ListenMessager()

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
