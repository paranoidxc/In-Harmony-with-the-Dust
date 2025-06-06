package main

import (
	"bigwork/client/test"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	c := test.NewClient()
	c.InputHandlerRegister()
	c.MessageHandlerRegister()
	c.Run()
	select {}
}
