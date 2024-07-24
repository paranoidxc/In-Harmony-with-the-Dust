package main

import (
	"bigwork/world"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	world.MM = world.NewMgrMgr()
	go world.MM.Run()
	select {}
}
