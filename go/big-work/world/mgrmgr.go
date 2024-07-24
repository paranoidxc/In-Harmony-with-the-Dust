package world

import (
	"bigwork/manager"
	"bigwork/network"
	"bigwork/network/protocol/gen/messageId"
	"log"
)

type MgrMgr struct {
	Pm       *manager.PlayerMgr
	Server   *network.Server
	Handlers map[messageId.MessageId]func(message *network.SessionPacket)
}

func NewMgrMgr() *MgrMgr {
	m := &MgrMgr{Pm: manager.NewPlayerMgr()}
	m.Server = network.NewServer(":8023", "tcp6")
	m.Server.OnSessionPacket = m.OnSessionPacket
	m.Handlers = make(map[messageId.MessageId]func(message *network.SessionPacket))

	return m
}

var MM *MgrMgr

func (mm *MgrMgr) Run() {
	mm.HandlerRegister()
	go mm.Server.Run()
	go mm.Pm.Run()
}

func (mm *MgrMgr) OnSessionPacket(packet *network.SessionPacket) {
	log.Println("MgrMgr OnSessionPacket", *packet)
	log.Println("Packet", packet.Msg.ID)
	if handler, ok := mm.Handlers[messageId.MessageId(packet.Msg.ID)]; ok {
		log.Println("MgrMgr mm Handlers")
		handler(packet)
		return
	}
	if p := mm.Pm.GetPlayer(packet.Sess.UId); p != nil {
		log.Println(" p.HandlerParamCh <- packet.Msg ")
		p.HandlerParamCh <- packet.Msg
	}
}
