package player

import "bigwork/network/protocol/gen/messageId"

func (p *Player) HandlerRegister() {
	p.handles[messageId.MessageId_CSAddFriend] = p.AddFriend
	p.handles[messageId.MessageId_CSDelFriend] = p.DelFriend
	p.handles[messageId.MessageId_CSSendChatMsg] = p.ResolveChatMsg
}
