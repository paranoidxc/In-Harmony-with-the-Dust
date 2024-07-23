package network

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type Client struct {
	Address string
	packer  *NormalPacker
	//Conn    net.Conn
}

func NewClient(address string) *Client {
	return &Client{
		Address: address,
		packer:  NewNormalPacker(binary.BigEndian),
	}
}

func (c *Client) Run() {
	conn, err := net.Dial("tcp6", c.Address)
	if err != nil {
		fmt.Println(err)
		return
	}

	//c.Conn = conn

	go c.Write(conn)
	go c.Read(conn)
}

func (c *Client) Write(conn net.Conn) {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
			c.send(conn, &Message{
				Id:   111,
				Data: []byte("Hello Fucking World"),
			})
		}
	}
}

func (c *Client) send(conn net.Conn, message *Message) {
	/*
		err := conn.SetWriteDeadline(time.Now().Add(time.Second))
		if err != nil {
			fmt.Println(err)
			return
		}
	*/

	bytes, err := c.packer.Pack(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("client write ok", string(message.Data))
}

func (c *Client) Read(conn net.Conn) {
	for {
		message, err := c.packer.Unpack(conn)
		if _, ok := err.(net.Error); err != nil && ok {
			fmt.Println(err)
			continue
		}

		fmt.Println("client receive message:", string(message.Data))
	}
}
