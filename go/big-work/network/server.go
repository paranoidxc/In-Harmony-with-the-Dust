package network

import (
	"log"
	"net"
)

type Server struct {
	tcpListener     net.Listener
	OnSessionPacker func(packet *SessionPacket)
}

func NewServer(address, network string) *Server {
	resolveTCPAddr, err := net.ResolveTCPAddr("tcp6", address)
	if err != nil {
		panic(err)
	}
	tcpListener, err := net.ListenTCP("tcp6", resolveTCPAddr)
	if err != nil {
		panic(err)
	}

	s := &Server{}
	s.tcpListener = tcpListener

	return s
}

func (s *Server) Run() {
	for {
		conn, err := s.tcpListener.Accept()
		log.Println("server accepting new client")
		if err != nil {
			continue
		}

		go func() {
			newSession := NewSession(conn)
			SessionMgrInstance.AddSession(newSession)
			newSession.Run()
			SessionMgrInstance.DelSession(newSession.UId)
		}()
	}
}
