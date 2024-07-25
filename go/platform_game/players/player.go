package players

import (
	"github.com/veandco/go-sdl2/sdl"
	"platform_game/colors"
	"platform_game/inputs"
	"platform_game/objects"
)

const (
	SPEED      = 8.0
	JUMP_FORCE = -16.0
)

type Player struct {
	*objects.Object
	Direction sdl.FPoint
	Speed     float32
	OnGround  bool
}

func New[size objects.Size](w, h size, x, y float32) *Player {
	player := new(Player)
	player.Object = objects.New(w, h, x, y)
	player.Speed = SPEED
	player.Color = colors.Red()
	return player
}

func (p *Player) Update() {

	if inputs.KeyPresssed(sdl.SCANCODE_RIGHT) {
		p.Direction.X = 1
	} else if inputs.KeyPresssed(sdl.SCANCODE_LEFT) {
		p.Direction.X = -1
	} else {
		p.Direction.X = 0
	}
	if inputs.KeyPresssed(sdl.SCANCODE_SPACE) && p.OnGround {
		p.Direction.Y = JUMP_FORCE
	}

	p.Position.X += p.Direction.X * p.Speed
	//log.Println("X", p.Position.X)
}
