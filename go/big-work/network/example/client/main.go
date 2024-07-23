package main

import (
	"bigwork/network"
	"fmt"
)

func main() {
	client := network.NewClient(":8088")
	client.Run()

	fmt.Println("client running")
	select {}
}
