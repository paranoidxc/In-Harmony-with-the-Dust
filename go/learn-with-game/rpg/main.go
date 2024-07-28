package main

import (
	"learn-with-game/rpg/game"
	"learn-with-game/rpg/ui2d"
)

func main() {
	//level := LoadLevelFromFile("game/maps/level1.map")
	num := 1
	game := game.NewGame(num)
	go func() {
		game.Run()
	}()
	ui := ui2d.NewUI(game.InputChan, game.LevelChans[0])
	ui.Run()
	/*
		time.Sleep(2 * time.Second)
		for i := 0; i < num; i++ {
			go func(i int) {
				//runtime.LockOSThread()
				ui := ui2d.NewUI(game.InputChan, game.LevelChans[i])
				ui.Run()
			}(i)
		}
	*/
}
