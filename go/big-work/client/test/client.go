package test

import (
	"bigwork/network"
	"encoding/json"
	"log"
)

type Client struct {
	cli             *network.Client
	inputHandlers   map[string]InputHandler
	messageHandlers map[uint64]MessageHandler
	chInput         chan *InputParam
	console         *ClientConsole
}

func NewClient() *Client {
	c := &Client{
		cli:             network.NewClient(":8023"),
		inputHandlers:   map[string]InputHandler{},
		messageHandlers: map[uint64]MessageHandler{},
		console:         NewClientConsole(),
	}
	c.cli.OnMessage = c.OnMessage
	c.chInput = make(chan *InputParam, 1)
	c.console.chInput = c.chInput

	c.MessageHandlerRegister()
	c.InputHandlerRegister()

	return c
}

func (c *Client) Run() {
	go func() {
		for {
			select {
			case input := <-c.chInput:
				log.Printf("Cmd:%s, prarms:%v <<< \t \n", input.Command, input.Param)
				bytes, err := json.Marshal(input.Param)
				if err == nil {
					c.cli.ChMsg <- &network.Message{
						ID:   111,
						Data: bytes,
					}
				}
			}
		}
	}()
	go c.console.Run()
	go c.cli.Run()
}

func (c *Client) OnMessage(packet *network.ClientPacket) {
	log.Println("test client OnMessage", *packet)
	if handler, ok := c.messageHandlers[packet.Msg.ID]; ok {
		handler(packet)
	} else {
		log.Println("handler not found")
	}
}
