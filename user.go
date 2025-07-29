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
}

// Create a new user
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// Start listening for messages
	go user.ListenMessage()

	return user
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
