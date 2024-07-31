package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
	"learn-with-game/rpg/game"
)

type mouseState struct {
	leftButton  bool
	rightButton bool
	pos         game.Pos
}

func getMouseState() *mouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()
	var result mouseState
	result.pos.X = int(mouseX)
	result.pos.Y = int(mouseY)
	result.leftButton = !(leftButton == 0)
	result.rightButton = !(rightButton == 0)
	return &result
}
