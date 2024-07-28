package game

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type Game struct {
	LevelChans []chan *Level
	InputChan  chan *Input
	Level      *Level
}

func NewGame(numWindows int, levelPath string) *Game {
	levelChans := make([]chan *Level, numWindows)
	for i := range levelChans {
		levelChans[i] = make(chan *Level)
	}
	inputChan := make(chan *Input)

	g := &Game{
		LevelChans: levelChans,
		InputChan:  inputChan,
		Level:      LoadLevelFromFile(levelPath),
	}

	return g
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	QuitGame
	CloseWindow
	Search
)

type Input struct {
	Typ          InputType
	LevelChannel chan *Level
}

type Tile struct {
	Rune        rune
	OverlayRune rune
	Visible     bool
	Seen        bool
}

const (
	StoreWall rune = '#'
	DirtFloor rune = '.'
	CloseDoor rune = '|'
	OpenDoor  rune = '/'
	UpStair   rune = 'u'
	DownStair rune = 'd'
	Blank     rune = 0
	Pending   rune = -1
)

type Pos struct {
	X, Y int
}

type Entity struct {
	Pos
	Rune rune
	Name string
}

type Character struct {
	Entity
	Hitpoints    int
	Strength     int
	Speed        float64
	ActionPoints float64
	SigntRange   int
}

type Player struct {
	Character
}

type Level struct {
	Map      [][]Tile
	Player   *Player
	Monsters map[Pos]*Monster
	Events   []string
	Debug    map[Pos]bool
	EventPos int
}

type Attackable interface {
	GetActionPoint() float64
	SetActionPoint(float64)
	GetHitpoints() int
	SetHitpoints(int)
	GetAttackPower() int
}

func (c *Character) GetActionPoint() float64 {
	return c.ActionPoints
}
func (c *Character) SetActionPoint(ap float64) {
	c.ActionPoints = ap
}
func (c *Character) GetHitpoints() int {
	return c.Hitpoints
}
func (c *Character) SetHitpoints(h int) {
	c.Hitpoints = h
}
func (c *Character) GetAttackPower() int {
	return c.Strength
}

func (level *Level) Attack(c1 *Character, c2 *Character) {
	c1.ActionPoints--
	c1AttackPower := c1.Strength
	c2.Hitpoints -= c1AttackPower

	if c2.Hitpoints > 0 {
		level.AddEvent(c1.Name + " Attacked " + c2.Name + " for " + strconv.Itoa(c1AttackPower))
	}
	//a1.SetActionPoint(a1.GetActionPoint() - 1)
	//a2.SetHitpoints(a2.GetHitpoints() - a1.GetAttackPower())
	/*
		if a2.GetHitpoints() > 0 {
			a2.SetActionPoint(a2.GetActionPoint() - 1)
			a1.SetHitpoints(a1.GetHitpoints() - a2.GetAttackPower())
		}
	*/
}

func PlayerAttackMonster(p *Player, m *Monster) {
	p.ActionPoints -= 1
	m.Hitpoints -= p.Strength

	if m.Hitpoints > 0 {
		m.ActionPoints -= 1
		p.Hitpoints -= m.Strength
	}
}
func MonsterAttackPlayer(m *Monster, p *Player) {
	m.ActionPoints -= 1
	p.Hitpoints -= m.Strength

	if p.Hitpoints > 0 {
		p.ActionPoints -= 1
		m.Hitpoints -= p.Strength
	}
}

/*
func Attack(c1 Character, c2 Character) {
	c1.ActionPoints -= 1
	c2.Hitpoints -= c1.Strength

	if c2.Hitpoints > 0 {
		c1.Hitpoints -= c2.Strength
		c2.ActionPoints -= 1
	}
}
*/

/*
old

	type priorityPos struct {
		Pos
		priority int
	}

type priorityArray []priorityPos

func (p priorityArray) Len() int           { return len(p) }
func (p priorityArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p priorityArray) Less(i, j int) bool { return p[i].priority < p[j].priority }
*/
func (level *Level) AddEvent(event string) {
	level.Events[level.EventPos] = event

	level.EventPos++
	if level.EventPos == len(level.Events) {
		level.EventPos = 0
	}
}

func LoadLevelFromFile(filename string) *Level {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	levelLines := make([]string, 0)
	longestRow := 0
	index := 0
	for scanner.Scan() {
		levelLines = append(levelLines, scanner.Text())
		if len(levelLines[index]) > longestRow {
			longestRow = len(levelLines[index])
		}
		index++
	}

	level := &Level{}
	level.Debug = make(map[Pos]bool)
	level.Events = make([]string, 10)
	level.Player = &Player{}
	level.Player.Strength = 5
	level.Player.Hitpoints = 120
	level.Player.Name = "GoMan"
	level.Player.Rune = '@'
	level.Player.Speed = 1.0
	level.Player.ActionPoints = 0
	level.Player.SigntRange = 5

	level.Map = make([][]Tile, len(levelLines))
	level.Monsters = make(map[Pos]*Monster)

	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}

	for y := 0; y < len(level.Map); y++ {
		line := levelLines[y]
		for x, c := range line {
			var t Tile
			t.OverlayRune = Blank
			switch c {
			case ' ', '\t', '\n', '\r':
				t.Rune = Blank
			case '#':
				t.Rune = StoreWall
			case '|':
				t.OverlayRune = CloseDoor
				t.Rune = Pending
			case '/':
				t.OverlayRune = OpenDoor
				t.Rune = Pending
			case 'u':
				t.OverlayRune = UpStair
				t.Rune = Pending
			case 'd':
				t.OverlayRune = DownStair
				t.Rune = Pending
			case '.':
				t.Rune = DirtFloor
			case '@':
				level.Player.X = x
				level.Player.Y = y
				t.Rune = Pending
			case 'R':
				level.Monsters[Pos{x, y}] = NewRat(Pos{x, y})
				t.Rune = Pending
			case 'S':
				level.Monsters[Pos{x, y}] = NewSpider(Pos{x, y})
				t.Rune = Pending
			default:
				//panic("Invalid character in map")
			}
			level.Map[y][x] = t
		}
	}

	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune == Pending {
				level.Map[y][x].Rune = level.bfsFloor(Pos{x, y})
				/*
					SearchLoop:
						for searchX := x - 1; searchX <= x+1; searchX++ {
							for searchY := y - 1; searchY <= y+1; searchY++ {
								searchTile := level.Map[searchY][searchX]
								switch searchTile {
								case DirtFloor:
									level.Map[y][x] = DirtFloor
									break SearchLoop
								}
							}
						}
				*/
			}
		}
	}

	level.lineOfSignt()
	return level
}

func inRange(level *Level, pos Pos) bool {
	return pos.X < len(level.Map[0]) && pos.Y < len(level.Map) && pos.X >= 0 && pos.Y >= 0
}

func canSeeThrough(level *Level, pos Pos) bool {
	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch t.Rune {
		case StoreWall, Blank:
			return false
		default:
			return true
		}

		switch t.OverlayRune {
		case CloseDoor:
			return false
		default:
			return true
		}
	}
	return false
}

func canWalk(level *Level, pos Pos) bool {
	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch t.Rune {
		case StoreWall, Blank:
			return false
		}
		switch t.OverlayRune {
		case CloseDoor:
			return false
		}

		_, exists := level.Monsters[pos]
		if exists {
			return false
		}

		return true
	}

	return false
}

func checkDoor(level *Level, pos Pos) {
	t := level.Map[pos.Y][pos.X]
	if t.OverlayRune == CloseDoor {
		level.Map[pos.Y][pos.X].OverlayRune = OpenDoor
		level.lineOfSignt()
	}
}

func (player *Player) Move(to Pos, level *Level) {
	player.Pos = to
	for y, row := range level.Map {
		for x, _ := range row {
			level.Map[y][x].Visible = false
		}
	}
	//pPos := player.Pos
	//level.bresenham(pPos, Pos{pPos.X, pPos.Y - player.SigntRange})
	level.lineOfSignt()
}

func (level *Level) resolveMovement(pos Pos) {
	monster, exists := level.Monsters[pos]
	if exists {
		level.Attack(&level.Player.Character, &monster.Character)
		if monster.Hitpoints <= 0 {
			delete(level.Monsters, monster.Pos)
		}

		if level.Player.Hitpoints <= 0 {
			panic("ded")
		}
	} else if canWalk(level, pos) {
		level.Player.Move(pos, level)
	} else {
		checkDoor(level, pos)
	}
}

func (game *Game) handleInput(input *Input) {
	level := game.Level
	p := level.Player

	switch input.Typ {
	case Up:
		newPos := Pos{p.X, p.Y - 1}
		level.resolveMovement(newPos)
	case Down:
		newPos := Pos{p.X, p.Y + 1}
		level.resolveMovement(newPos)
	case Left:
		newPos := Pos{p.X - 1, p.Y}
		level.resolveMovement(newPos)
	case Right:
		newPos := Pos{p.X + 1, p.Y}
		level.resolveMovement(newPos)
	case Search:
		//bfs(ui, Level, p.Pos)
		level.astar(level.Player.Pos, Pos{level.Player.X + 2, level.Player.Y + 1})
	case CloseWindow:
		close(input.LevelChannel)
		chanIndex := 0
		for i, c := range game.LevelChans {
			if c == input.LevelChannel {
				chanIndex = i
				break
			}
		}
		game.LevelChans = append(game.LevelChans[:chanIndex], game.LevelChans[chanIndex+1:]...)
	}
}

func getNeighbors(level *Level, pos Pos) []Pos {
	neighbors := make([]Pos, 0, 4)
	left := Pos{pos.X - 1, pos.Y}
	right := Pos{pos.X + 1, pos.Y}
	up := Pos{pos.X, pos.Y - 1}
	down := Pos{pos.X, pos.Y + 1}

	if canWalk(level, right) {
		neighbors = append(neighbors, right)
	}
	if canWalk(level, left) {
		neighbors = append(neighbors, left)
	}
	if canWalk(level, up) {
		neighbors = append(neighbors, up)
	}
	if canWalk(level, down) {
		neighbors = append(neighbors, down)
	}

	return neighbors
}

func (level *Level) bfsFloor(start Pos) rune {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true

	//level.Debug = visited

	for len(frontier) > 0 {
		current := frontier[0]
		currentTile := level.Map[current.Y][current.X]
		switch currentTile.Rune {
		case DirtFloor:
			return DirtFloor
			//return Tile{DirtFloor, Blank, false, false}
		default:
		}

		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}
	}

	return DirtFloor
}

func (level *Level) bfs(start Pos) {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true

	//level.Debug = visited

	for len(frontier) > 0 {
		current := frontier[0]
		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (level *Level) astar(start Pos, goal Pos) []Pos {
	frontier := make(pqueue, 0, 8)
	frontier = frontier.push(start, 1)
	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start
	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0

	//level.Debug = make(map[Pos]bool)

	current := Pos{}
	for len(frontier) > 0 {
		frontier, current = frontier.pop()
		if current == goal {
			path := make([]Pos, 0)
			p := current
			for p != start {
				path = append(path, p)
				p = cameFrom[p]
			}
			path = append(path, p)

			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			//for _, pos := range path {
			//level.Debug[pos] = true
			//time.Sleep(100 * time.Millisecond)
			//}
			return path
		}

		for _, next := range getNeighbors(level, current) {
			newCost := costSoFar[current] + 1
			_, exists := costSoFar[next]
			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
				frontier = frontier.push(next, priority)
				cameFrom[next] = current
			}
		}
	}

	return nil
}

func (game *Game) Run() {
	fmt.Println("Starting game...")
	//Level := LoadLevelFromFile("game/maps/level1.map")
	for _, lchan := range game.LevelChans {
		lchan <- game.Level
	}

	for input := range game.InputChan {
		if input != nil && input.Typ == QuitGame {
			fmt.Println("QUIT")
			return
		}
		//p := game.Level.Player.Pos
		//level.bresenham(p, Pos{p.X + 5, p.Y + 5})

		game.handleInput(input)
		for _, monster := range game.Level.Monsters {
			monster.Update(game.Level)
		}

		/*
			if len(game.LevelChans) == 0 {
				return
			}
		*/

		for _, lchan := range game.LevelChans {
			lchan <- game.Level
		}
	}
}
