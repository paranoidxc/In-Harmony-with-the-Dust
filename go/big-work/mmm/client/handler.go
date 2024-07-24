package main

import "bigwork/network"

type MessageHandler func(packet *network.ClientPacket)

type InputHandler func(param *InputParam)

func (c *Client) Login(param *InputParam) {

}

func (c *Client) OnLoginRsp(packet *network.ClientPacket) {

}

func (c *Client) AddFriend(param *InputParam) {

}

func (c *Client) OnAddFriendRsp(packet *network.ClientPacket) {

}

func (c *Client) DelFriend(param *InputParam) {

}

func (c *Client) OnDelFriendRsp(packet *network.ClientPacket) {

}

func (c *Client) SendChatMsg(param *InputParam) {

}

func (c *Client) OnSendChatMsgRsp(packet *network.ClientPacket) {

}
