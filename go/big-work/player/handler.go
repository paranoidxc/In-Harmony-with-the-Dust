package player

import "bigwork/network"

type Handler func(packet *network.Message)

func (p *Player) AddFriend(packet *network.Message) {
}

func (p *Player) DelFriend(packet *network.Message) {
}

func (p *Player) ResolveChatMsg(packet *network.Message) {
}
