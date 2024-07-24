package player

import "bigwork/network"

type Player struct {
	UId           uint64
	FriendList    []uint64
	HandleParamCh chan *network.Message
	handles       map[uint64]Handler
	session       *network.Session
}

func NewPlayer() *Player {
	p := &Player{
		UId:        0,
		FriendList: make([]uint64, 100),
		handles:    make(map[uint64]Handler),
	}

	p.HandlerRegister()
	return p
}

func (p *Player) Run() {
	for {
		select {
		case HandleParam := <-p.HandleParamCh:
			if fn, ok := p.handles[HandleParam.ID]; ok {
				fn(HandleParam.Data)
			}
		}
	}
}
