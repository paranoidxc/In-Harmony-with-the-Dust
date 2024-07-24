package player

import (
	"bigwork/network"
	"bigwork/network/protocol/gen/messageId"
)

type Player struct {
	UId            uint64
	FriendList     []uint64
	HandlerParamCh chan *network.Message
	handles        map[messageId.MessageId]Handler
	Session        *network.Session
}

func NewPlayer() *Player {
	p := &Player{
		UId:        0,
		FriendList: make([]uint64, 100),
		handles:    make(map[messageId.MessageId]Handler),
	}

	p.HandlerRegister()
	return p
}

func (p *Player) Run() {
	for {
		select {
		case HandlerParam := <-p.HandlerParamCh:
			if fn, ok := p.handles[messageId.MessageId(HandlerParam.ID)]; ok {
				fn(HandlerParam)
			}
		}
	}
}
