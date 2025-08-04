package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// Create a new user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	// Start listening for messages
	go user.ListenMessage()

	return user
}

func (u *User) Online() {

	// Add the user to the online user map
	u.server.maplock.Lock()
	u.server.OnlineUserMap[u.Name] = u
	u.server.maplock.Unlock()

	// Broadcast that the user has come online
	u.server.BroadCast(u, "\033[32mhas come online\033[0m")
}

func (u *User) Offline() {
	// Remove the user from the online user map
	u.server.maplock.Lock()
	delete(u.server.OnlineUserMap, u.Name)
	u.server.maplock.Unlock()

	// Broadcast that the user has come online
	u.server.BroadCast(u, "\033[36mhas gone offline\033[0m")
}

func (u *User) DoMessage(msg string) {
	u.server.BroadCast(u, msg)
}
	
// listen the user's channel, if there is a message, send it to the user
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		if _, err := u.conn.Write([]byte(msg + "\n")); err != nil {
			fmt.Println("Error writing to user:", u.Name, err)
			u.conn.Close()
			return
		}
	}
}
