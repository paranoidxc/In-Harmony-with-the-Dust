package player

func (p *Player) HandlerRegister() {
	p.handles[111] = p.AddFriend
	p.handles[222] = p.DelFriend
	p.handles[333] = p.ResolveChatMsg
}
