package game

import "fmt"

type Monster struct {
	Character
}

func NewRat(p Pos) *Monster {
	//return &Monster{X: pos.X, Y: pos.Y, Rune: 'R', Name: "Rat", Hitpoints: 5, Strength: 5, Speed: 2.0, ActionPoints: 0.0}
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'R'
	monster.Name = "Rat"
	monster.Hitpoints = 5
	monster.Strength = 5
	monster.Speed = 1.5
	monster.ActionPoints = 0.0
	return monster
}

func NewSpider(p Pos) *Monster {
	/*
		return &Monster{
			Pos:          pos,
			Rune:         'S',
			Name:         "Spider",
			Hitpoints:    10,
			Strength:     10,
			Speed:        1.0,
			ActionPoints: 0.0,
		}
	*/
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'S'
	monster.Name = "Rat"
	monster.Hitpoints = 10
	monster.Strength = 10
	monster.Speed = 1.0
	monster.ActionPoints = 0.0
	return monster
}

func (m *Monster) Update(level *Level) {
	m.ActionPoints += m.Speed
	playerPos := level.Player.Pos

	apInt := int(m.ActionPoints)
	positions := level.astar(m.Pos, playerPos)
	moveIndex := 1
	for i := 0; i < apInt; i++ {
		if moveIndex < len(positions) {
			m.move(positions[moveIndex], level)
			moveIndex++
			m.ActionPoints--
		}
	}
}

func (m *Monster) move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
	} else {
		Attack(m, level.Player)
		if m.Hitpoints <= 0 {
			delete(level.Monsters, m.Pos)
		}

		if level.Player.Hitpoints <= 0 {
			fmt.Println("YOU DIED")
			panic("YOU DIED")
		}
	}
}
