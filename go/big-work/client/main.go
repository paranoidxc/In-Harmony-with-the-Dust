package main

import "bigwork/client/test"

func main() {
	c := test.NewClient()
	c.Run()
	select {}
}
