package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
			// Use select to prevent writing to a closed channel.
			select {
			case user.C <- msg:
				// Message sent successfully
			default:
				// Channel is full or closed, skip this user
				fmt.Printf("Failed to send message to user %s, channel may be closed\n", user.Name)
			}
		}
		s.maplock.RUnlock()
	}
}

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]\033[35m" + user.Name + "\033[0m:" + msg

	s.MessageChan <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	// Handle the connection
	fmt.Printf("\033[32mClient connection established successfully from:\033[0m (%s)%s\n", conn.RemoteAddr().Network(), conn.RemoteAddr().String())

	user := NewUser(conn, s)

	user.Online()

	// Notify the user that they are online
	isLive := make(chan bool)

	// Accept messages from the user
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
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

			// any message received, reset the isLive channel
			isLive <- true
		}
	}()

	// Block the handler
	for {
		select {
		case <-isLive:
			// User is active, reset the timer
		case <-time.After(30 * time.Minute):
			// timeout, close the connection
			user.SendMsg("\033[33mYou have been inactive for too long, disconnecting...\033[0m\n")

			// First let the user go offline (remove from online list)
			user.Offline()

			// Close user connection and resources
			close(user.C) // Close the user's message channel
			conn.Close()  // Close the connection
			fmt.Printf("[(%s)%s]\033[35m%s:\033[34mhas been disconnected due to inactivity\033[0m\n", conn.RemoteAddr().Network(), user.Addr, user.Name)

			return
		}
	}
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

	// Add server console input listener
	quitChan := make(chan bool)
	go func() {
		for {
			var input string
			fmt.Scanln(&input)
			if input == "QUIT" || input == "quit" {
				fmt.Println("\033[31mServer is shutting down...\033[0m")
				quitChan <- true
				return
			}
		}
	}()

	for {
		select {
		case <-quitChan:
			// Received exit signal, shutting down the server.
			fmt.Println("\033[33mClosing all connections...\033[0m")

			// Close all user connections
			s.maplock.Lock()
			for _, user := range s.OnlineUserMap {
				user.SendMsg("\033[31mServer is shutting down. Connection will be closed.\033[0m\n")
				close(user.C)     // First close the channel
				user.conn.Close() // Then close the connection
			}
			// Clear the online user map
			s.OnlineUserMap = make(map[string]*User)
			s.maplock.Unlock()

			fmt.Println("\033[32mServer stopped successfully.\033[0m")
			return

		default:
			// Set the listener to non-blocking mode and check for new connections
			listener.(*net.TCPListener).SetDeadline(time.Now().Add(100 * time.Millisecond))
			conn, err := listener.Accept()

			if err != nil {
				// Check if it's a timeout error (normal case)
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // Timeout is normal, continue the loop
				}
				fmt.Println("listener.Accept error:", err)
				continue
			}

			// handle connections
			go s.Handler(conn)
		}
	}
}
