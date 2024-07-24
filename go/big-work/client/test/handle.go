package test

import (
	"bigwork/network"
	"bigwork/network/protocol/gen/player"
	"google.golang.org/protobuf/proto"
	"log"
	"strconv"
)

type MessageHandler func(packet *network.ClientPacket)
type InputHandler func(param *InputParam)

// CreatePlayer 创建角色
func (c *Client) CreatePlayer(param *InputParam) {
	id := c.GetMessageIdByCmd(param.Command)
	if len(param.Param) != 2 {
		return
	}

	msg := &player.CSCreateUser{
		UserName: param.Param[0],
		Password: param.Param[1],
	}
	c.Transport(id, msg)
}

func (c *Client) OnCreatePlayerRsp(packet *network.ClientPacket) {
	log.Println("恭喜你创建角色成功")
}

func (c *Client) Login(param *InputParam) {
	log.Printf("Login input Handler print")
	log.Println(param.Command)
	log.Println(param.Param)

	id := c.GetMessageIdByCmd(param.Command)
	if len(param.Param) != 2 {
		return
	}

	msg := &player.CSLogin{
		UserName: param.Param[0],
		Password: param.Param[1],
	}
	c.Transport(id, msg)
}

func (c *Client) OnLoginRsp(packet *network.ClientPacket) {
	log.Println("OnLoginRsp")
	rsp := &player.SCLogin{}
	err := proto.Unmarshal(packet.Msg.Data, rsp)
	if err != nil {
		return
	}
	log.Println("LOGIN SUC")
}

func (c *Client) AddFriend(param *InputParam) {
	id := c.GetMessageIdByCmd(param.Command)
	// xxx || 空串
	if len(param.Param) != 1 || len(param.Param[0]) == 0 {
		return
	}

	uid, err := strconv.ParseUint(param.Param[0], 10, 64)
	if err != nil {
		return
	}

	msg := &player.CSAddFriend{
		UId: uid,
	}

	c.Transport(id, msg)
}

func (c *Client) OnAddFriendRsp(packet *network.ClientPacket) {

}

func (c *Client) DelFriend(param *InputParam) {
	id := c.GetMessageIdByCmd(param.Command)
	// xxx || 空串
	if len(param.Param) != 1 || len(param.Param[0]) == 0 {
		return
	}

	uid, err := strconv.ParseUint(param.Param[0], 10, 64)
	if err != nil {
		return
	}

	msg := &player.CSDelFriend{
		UId: uid,
	}

	c.Transport(id, msg)
}

func (c *Client) OnDelFriendRsp(packet *network.ClientPacket) {
	log.Println("you have del friend success")
}

func (c *Client) SendChatMsg(param *InputParam) {
	log.Println("SendChatMsg")
	id := c.GetMessageIdByCmd(param.Command)
	if len(param.Param) != 3 {
		log.Println("len(param.Param)  err!")
		return
	}

	uid, err := strconv.ParseUint(param.Param[0], 10, 64)
	if err != nil {
		log.Println("strconv.ParseUint 1 err!", err)
		return
	}

	log.Println("param", param.Param)
	parseInt32, err := strconv.ParseUint(param.Param[2], 10, 32)
	if err != nil {
		log.Println("strconv.ParseUint 2  err!", err)
		return
	}

	msg := &player.CSSendChatMsg{
		UId: uid,
		Msg: &player.ChatMessage{
			Content: param.Param[1],
			Extra:   nil,
		},
		Category: int32(parseInt32),
	}

	c.Transport(id, msg)
}

func (c *Client) OnSendChatMsgRsp(packet *network.ClientPacket) {
	log.Println("send  chat message success")
}
