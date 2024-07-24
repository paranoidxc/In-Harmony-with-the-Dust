package network

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

type Client struct {
	Address   string
	packer    IPacker
	ChMsg     chan *Message
	OnMessage func(message *ClientPacket)
}

func NewClient(address string) *Client {
	return &Client{
		Address:   address,
		packer:    &NormalPacker{ByteOrder: binary.BigEndian},
		ChMsg:     make(chan *Message, 1),
		OnMessage: clientHandleMsg,
	}
}

func (c *Client) Run() {
	conn, err := net.Dial("tcp6", c.Address)
	if err != nil {
		fmt.Println(err)
		return
	}

	//c.Conn = conn

	go c.Read(conn)
	go c.Write(conn)
}

func (c *Client) Write(conn net.Conn) {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
			log.Println("tick msg")
			c.ChMsg <- &Message{
				ID:   111,
				Data: []byte("Hello Fucking World"),
			}
		case msg := <-c.ChMsg:
			log.Println("send msg")
			c.send(conn, msg)
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
		c.OnMessage(&ClientPacket{
			Msg:  message,
			Conn: conn,
		})
		fmt.Println("client read msg:", string(message.Data))
	}
}

func clientHandleMsg(packet *ClientPacket) {
	log.Printf("clientHandleMsg: %+v\n", packet)
	log.Println(packet)
}
