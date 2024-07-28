package game

import "fmt"

type Monster struct {
	Character
}

func NewRat(p Pos) *Monster {
	//return &Monster{Character: Character{}}
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'R'
	monster.Name = "Rat"
	monster.Hitpoints = 20
	monster.Strength = 1
	monster.Speed = 1.5
	monster.ActionPoints = 0.0
	return monster
}

func NewSpider(p Pos) *Monster {
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'S'
	monster.Name = "Rat"
	monster.Hitpoints = 40
	monster.Strength = 1
	monster.Speed = 1.0
	monster.ActionPoints = 0.0
	return monster
}

func (m *Monster) Update(level *Level) {
	m.ActionPoints += m.Speed
	playerPos := level.Player.Pos

	apInt := int(m.ActionPoints)
	positions := level.astar(m.Pos, playerPos)

	// no path to player
	if len(positions) == 0 {
		m.Pass()
		return
	}
	moveIndex := 1
	for i := 0; i < apInt; i++ {
		if moveIndex < len(positions) {
			m.move(positions[moveIndex], level)
			moveIndex++
			m.ActionPoints--
		}
	}
}

func (m *Monster) Pass() {
	m.ActionPoints -= m.Speed
}

func (m *Monster) move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
		return
	}

	if to == level.Player.Pos {
		level.AddEvent(m.Name + " Attack Player")
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
