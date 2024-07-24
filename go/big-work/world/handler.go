package world

import (
	"bigwork/network"
	"bigwork/network/protocol/gen/messageId"
	"bigwork/network/protocol/gen/player"
	logicPlayer "bigwork/player"
	"google.golang.org/protobuf/proto"
	"log"
)

func (mm *MgrMgr) CreatePlayer(message *network.SessionPacket) {
	log.Println("[MgrMgr.CreatePlayer Call]")
	msg := &player.CSCreateUser{}
	err := proto.Unmarshal(message.Msg.Data, msg)
	if err != nil {
		return
	}
	log.Println("[MgrMgr.CreatePlayer]", msg)
	mm.SendMsg(uint64(messageId.MessageId_SCCreatePlayer), &player.SCCreateUser{}, message.Sess)
}

func (mm *MgrMgr) UserLogin(message *network.SessionPacket) {
	log.Println("[MgrMgr.UserLogin Call]")
	msg := &player.CSLogin{}
	err := proto.Unmarshal(message.Msg.Data, msg)
	if err != nil {
		return
	}
	newPlayer := logicPlayer.NewPlayer()
	newPlayer.UId = 111
	//newPlayer.UId = uint64(time.Now().Unix())
	newPlayer.HandlerParamCh = message.Sess.WriteCh
	message.Sess.IsPlayerOnline = true
	message.Sess.UId = newPlayer.UId
	newPlayer.Session = message.Sess
	mm.Pm.Add(newPlayer)
}

func (mm *MgrMgr) SendMsg(id uint64, message proto.Message, session *network.Session) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		return
	}
	rsp := &network.Message{
		ID:   id,
		Data: bytes,
	}
	session.SendMsg(rsp)
}
