package world

import "bigwork/network/protocol/gen/messageId"

func (mm *MgrMgr) HandlerRegister() {
	mm.Handlers[messageId.MessageId_CSCreatePlayer] = mm.CreatePlayer
	mm.Handlers[messageId.MessageId_CSLogin] = mm.UserLogin
}
