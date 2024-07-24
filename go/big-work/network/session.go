package network

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type Session struct {
	UId            int64
	conn           net.Conn
	IsClose        bool
	packer         IPacker
	WriteCh        chan *Message
	IsPlayerOnline bool
	MessageHandler func(packet *SessionPacket)
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		conn:           conn,
		packer:         &NormalPacker{ByteOrder: binary.BigEndian},
		WriteCh:        make(chan *Message, 1),
		MessageHandler: serverHandleMsg,
	}
}

func (s *Session) Run() {
	go s.Read()
	go s.Write()
}

func (s *Session) Read() {
	//fmt.Println("session read")
	/*
		err := s.conn.SetReadDeadline(time.Now().Add(time.Second))
		if err != nil {
			fmt.Println("SetReadDeadline", err)
		}
	*/
	for {
		//fmt.Println("server before read")
		message, err := s.packer.Unpack(s.conn)
		if _, ok := err.(net.Error); ok {
			fmt.Println(err)
			continue
		}
		if message == nil {
			continue
		}
		//fmt.Printf("server receive message: %+v", message)
		fmt.Println("server receive message:", string(message.Data))
		/*
			s.MessageHandler(&SessionPacket{
				Msg:  message,
				Sess: s,
			})
		*/
		s.WriteCh <- &Message{
			ID:   111,
			Data: []byte("Hi"),
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
		case msg := <-s.WriteCh:
			//fmt.Println("channel", *msg)
			s.send(msg)
		}
	}
}

func (s *Session) send(message *Message) {
	/*
		err := s.Conn.SetWriteDeadline(time.Now().Add(time.Second))
		if err != nil {
			fmt.Println(err)
			return
		}
	*/

	bytes, err := s.packer.Pack(message)
	if err != nil {
		return
	}
	_, err = s.conn.Write(bytes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("server write message ok:", message)
}

func serverHandleMsg(packet *SessionPacket) {
	log.Printf("serverHandleMsg: %+v\n", packet)
	log.Println(packet)
}
