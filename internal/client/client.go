package client

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,

		flag: 666, // default value for flag
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

// deal server responses
func (c *Client) DealResponse() {
	// Once client.conn has data, copy it to stdout, block listening forever
	io.Copy(os.Stdout, c.conn)
}

func (c *Client) menu() bool {
	var flag int

	fmt.Println("\033[34m1. Public chat\033[0m")
	fmt.Println("\033[34m2. Private chat\033[0m")
	fmt.Println("\033[34m3. Update user name\033[0m")
	fmt.Println("\033[34m0. Exit\033[0m")

	_, err := fmt.Scanln(&flag)
	if err != nil {
		fmt.Println("\033[31m>>>>>Invalid input, please try again<<<<<\033[0m")
		return false
	}

	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println("\033[33m>>>>>Invalid input, please try again<<<<<\033[0m")
		return false
	}
}

func (c *Client) PublicChat() {
	var chatMsg string

	fmt.Println("\033[34m[Input exit to exit](Public chat)\033[0m")
	_, err := fmt.Scanln(&chatMsg)
	if err != nil {
		return
	}

	for chatMsg != "exit" {

		// message cannot be empty
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err = c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("\033[31mconn.Write err:\033[0m", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("\033[34m[Input exit to exit](Public chat)\033[0m")
		_, err = fmt.Scanln(&chatMsg)
		if err != nil {
			return
		}

	}
}

func (c *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("\033[31mconn.Write err:\033[0m", err)
		return
	}
}

func (c *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	c.SelectUser()
	fmt.Println("\033[34mPlease input the user name you want to chat with:\033[0m")
	_, err := fmt.Scanln(&remoteName)
	if err != nil {
		return
	}

	for remoteName != "exit" {
		fmt.Println("\033[34m[Input exit to exit](Private chat)\033[0m")
		_, err = fmt.Scanln(&chatMsg)
		if err != nil {
			return
		}

		for chatMsg != "exit" {
			// message cannot be empty
			if len(chatMsg) != 0 {
				sendMsg := "tell|" + remoteName + "|" + chatMsg + "\n\n"
				_, err = c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("\033[31mconn.Write err:\033[0m", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println("\033[34m[Input exit to exit](Public chat)\033[0m")
			_, err = fmt.Scanln(&chatMsg)
			if err != nil {
				return
			}
		}

		c.SelectUser()
		fmt.Println("\033[34mPlease input the user name you want to chat with:\033[0m")
		_, err = fmt.Scanln(&remoteName)
		if err != nil {
			return
		}
	}

}

func (c *Client) UpdateName() bool {
	fmt.Println("\033[34m>>>>>Please input user name:\033[0m")
	_, err := fmt.Scanln(&c.Name)
	if err != nil {
		return false
	}

	sendMsg := "rename|" + c.Name + "\n"
	_, err = c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("\033[31mconn.Write err:\033[0m", err)
		return false
	}

	return true
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {
		}

		switch c.flag {
		case 1:
			fmt.Println("\033[34mPublic chat selected\033[0m")
			c.PublicChat()
		case 2:
			fmt.Println("\033[34mPrivate chat selected\033[0m")
			c.PrivateChat()
		case 3:
			fmt.Println("\033[34mUpdate user name selected\033[0m")
			c.UpdateName()
		case 0:
			fmt.Println("Exiting...")
		}
	}
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
