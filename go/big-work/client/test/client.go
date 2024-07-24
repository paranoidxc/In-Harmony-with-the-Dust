package test

import (
	"bigwork/network"
	"bigwork/network/protocol/gen/messageId"
	"log"
)

type Client struct {
	cli             *network.Client
	inputHandlers   map[string]InputHandler
	messageHandlers map[messageId.MessageId]MessageHandler
	chInput         chan *InputParam
	console         *ClientConsole
}

func NewClient() *Client {
	c := &Client{
		cli:             network.NewClient(":8023"),
		inputHandlers:   map[string]InputHandler{},
		messageHandlers: map[messageId.MessageId]MessageHandler{},
		console:         NewClientConsole(),
	}
	c.cli.OnMessage = c.OnMessage
	c.chInput = make(chan *InputParam, 1)
	c.console.chInput = c.chInput

	return c
}

func (c *Client) Run() {
	go func() {
		for {
			select {
			case input := <-c.chInput:
				log.Printf("Cmd:%s, prarms:%v <<< \t \n", input.Command, input.Param)
				inputHandler := c.inputHandlers[input.Command]
				if inputHandler != nil {
					log.Println("inputHandler call")
					inputHandler(input)
				} else {
					log.Println("inputHandler nil")
				}

				/*
					bytes, err := json.Marshal(input.Param)
					if err == nil {
						c.cli.ChMsg <- &network.Message{
							ID:   111,
							Data: bytes,
						}
					}
				*/
			}
		}
	}()
	go c.console.Run()
	go c.cli.Run()
}

func (c *Client) OnMessage(packet *network.ClientPacket) {
	log.Println("test client OnMessage", *packet)
	log.Println("Msg.ID", packet.Msg.ID)
	if handler, ok := c.messageHandlers[messageId.MessageId(packet.Msg.ID)]; ok {
		handler(packet)
	} else {
		log.Println("handler not found")
	}
}
