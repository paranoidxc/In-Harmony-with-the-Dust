package main

import (
	"bigwork/network"
	"log"
)

func main() {
	server := network.NewServer(":8023", "tcp6")
	log.Println("server start:8023")
	server.Run()
	select {}
}
