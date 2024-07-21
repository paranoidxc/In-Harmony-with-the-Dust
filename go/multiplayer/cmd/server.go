package main

import (
	"fmt"
	"log"
	"multiplayer-game-with-go/pkg/backend"
	"multiplayer-game-with-go/pkg/server"
	"multiplayer-game-with-go/proto"
	"net"

	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	port := 8888

	log.Printf("listening on port %d", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to Listen: %v", err)
	}
	fmt.Println(listener)

	game := backend.NewGame()
	game.Start()

	s := grpc.NewServer()
	server := server.NewGameServer(game, "")
	proto.RegisterGameServer(s, server)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
