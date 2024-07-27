package game

type Monster struct {
	Pos
	Rune      rune
	Name      string
	Hitpoints int
	Strength  int
	Speed     float64
}

func NewRat(pos Pos) *Monster {
	return &Monster{
		Pos:       pos,
		Rune:      'R',
		Name:      "Rat",
		Hitpoints: 5,
		Strength:  5,
		Speed:     2.0,
	}
}

func NewSpider(pos Pos) *Monster {
	return &Monster{
		Pos:       pos,
		Rune:      'S',
		Name:      "Spider",
		Hitpoints: 10,
		Strength:  10,
		Speed:     1.0,
	}
}

func (m *Monster) Update(level *Level) {
	playerPos := level.Player.Pos
	positions := level.astar(m.Pos, playerPos)

	// must len > 1 ; 1 is the monster's current position
	if len(positions) > 1 {
		m.move(positions[1], level)
	}
}

func (m *Monster) move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
	}
}
