package main

import (
	"learn-with-game/rpg/game"
	"learn-with-game/rpg/ui2d"
)

func main() {
	ui := &ui2d.UI2d{}
	game.Run(ui)
}
