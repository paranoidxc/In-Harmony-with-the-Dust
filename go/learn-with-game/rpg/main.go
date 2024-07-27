package main

import (
	"learn-with-game/rpg/game"
	"learn-with-game/rpg/ui2d"
)

func main() {
	num := 1
	//level := LoadLevelFromFile("game/maps/level1.map")
	game := game.NewGame(num, "game/maps/level1.map")
	go func() {
		game.Run()
	}()

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

	ui := ui2d.NewUI(game.InputChan, game.LevelChans[0])
	ui.Run()

	select {}
	/*
		for i := 0; i < num; i++ {
			go func(i int) {
				//runtime.LockOSThread()
				ui := ui2d.NewUI(game.InputChan, game.LevelChans[i])
				ui.Run()
			}(i)
		}
	*/
}
