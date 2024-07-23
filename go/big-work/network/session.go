package network

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Session struct {
	conn    net.Conn
	packer  *NormalPacker
	chWrite chan *Message
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		conn:    conn,
		packer:  NewNormalPacker(binary.BigEndian),
		chWrite: make(chan *Message, 1),
	}
}

func (s *Session) Run() {
	go s.Read()
	go s.Write()
}

func (s *Session) Read() {
	//fmt.Println("session read")
	/*8
	err := s.conn.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		fmt.Println("SetReadDeadline", err)
	}
	*/
	for {
		//fmt.Println("server before read")
		message, err := s.packer.Unpack(s.conn)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Printf("server receive message: %+v", message)
		fmt.Println("server receive message:", string(message.Data))
		s.chWrite <- &Message{
			Id:   999,
			Data: []byte("Hi" + string(message.Data)),
		}
	}
}

func (s *Session) Write() {
	//fmt.Println("session write")
	/*
		err := s.conn.SetWriteDeadline(time.Now().Add(time.Second))
		if err != nil {
			fmt.Println(err)
		}
	*/

	for {
		select {
		case msg := <-s.chWrite:
			//fmt.Println("channel", *msg)
			s.send(msg)
		}
	}
}

func (s *Session) send(message *Message) {
	bytes, err := s.packer.Pack(message)
	if err != nil {
		return
	}

	_, err = s.conn.Write(bytes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("server write mesage ok:", message)
}
