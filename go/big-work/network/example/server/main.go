package main

import "bigwork/network"

func main() {
	server := network.NewServer(":8088", "tcp6")
	server.Run()
	select {}
}
