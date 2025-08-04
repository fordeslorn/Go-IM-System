package main

import (
	"fmt"
	"net"
	"strings"
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

// Send a message to the user -> client
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

// handle user's message
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		// search the online user list
		u.server.maplock.RLock()
		for _, user := range u.server.OnlineUserMap {
			onlineMsg := "[" + user.Addr + "]\033[35m" + user.Name + "\033[0m:" + "\033[32monline\033[0m" + "\n"
			u.SendMsg(onlineMsg)
		}
		u.server.maplock.RUnlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// message format: rename|Jack
		newName := strings.Split(msg, "|")[1]

		// Judge if the new name already exists
		_, ok := u.server.OnlineUserMap[newName]
		if ok {
			u.SendMsg("\033[31mThe name already exists\033[0m\n")
		} else {
			u.server.maplock.Lock()
			delete(u.server.OnlineUserMap, u.Name)
			u.server.OnlineUserMap[newName] = u
			u.server.maplock.Unlock()

			u.Name = newName
			u.SendMsg("\033[36mYou have renamed yourself to:\033[35m" + newName + "\033[0m\n")
		}

	} else {
		u.server.BroadCast(u, msg)
	}
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
