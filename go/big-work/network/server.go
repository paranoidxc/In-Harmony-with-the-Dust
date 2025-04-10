package network

import (
	"log"
	"net"
)

type Server struct {
	tcpListener     net.Listener
	OnSessionPacket func(packet *SessionPacket)
}

func NewServer(address, network string) *Server {
	resolveTCPAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		panic(err)
	}
	tcpListener, err := net.ListenTCP(network, resolveTCPAddr)
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
			newSession.MessageHandler = s.OnSessionPacket
			SessionMgrInstance.AddSession(newSession)
			newSession.Run()
			SessionMgrInstance.DelSession(newSession.UId)
		}()
	}
}
