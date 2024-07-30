package game

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Game struct {
	LevelChans   []chan *Level
	InputChan    chan *Input
	Levels       map[string]*Level
	CurrentLevel *Level
}

func NewGame(numWindows int) *Game {
	levelChans := make([]chan *Level, numWindows)
	for i := range levelChans {
		levelChans[i] = make(chan *Level)
	}
	inputChan := make(chan *Input)

	levels := LoadLevels()
	g := &Game{
		LevelChans:   levelChans,
		InputChan:    inputChan,
		Levels:       levels,
		CurrentLevel: nil,
	}
	g.LoadWordFile()
	g.CurrentLevel.lineOfSight()

	return g
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	TakeItem
	QuitGame
	CloseWindow
	Search
	TakeAll
)

type Input struct {
	Typ          InputType
	Item         *Item
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

type LevelPos struct {
	*Level
	Pos
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
	Items        []*Item
}

type Player struct {
	Character
}

type GameEvent int

const (
	Move GameEvent = iota
	DoorOpen
	Attack
	Hit
	Portal
	PickUp
)

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

func (game *Game) LoadWordFile() {
	file, err := os.Open("game/maps/world.txt")
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	for rowIndex, row := range rows {
		if rowIndex == 0 {
			game.CurrentLevel = game.Levels[row[0]]
			continue
		}
		levelWithPortal := game.Levels[row[0]]
		if levelWithPortal == nil {
			fmt.Println("couldn't find level name in world file")
			panic(nil)
		}
		x, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			panic(err)
		}
		pos := Pos{int(x), int(y)}
		levelToTeleportTo := game.Levels[row[3]]

		x, err = strconv.ParseInt(row[4], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err = strconv.ParseInt(row[5], 10, 64)
		if err != nil {
			panic(err)
		}
		posToTeleportTo := Pos{int(x), int(y)}

		levelWithPortal.Portals[pos] = &LevelPos{levelToTeleportTo, posToTeleportTo}
	}
}

func LoadLevels() map[string]*Level {
	player := &Player{}
	player.Strength = 5
	player.Hitpoints = 120
	player.Name = "GoMan"
	player.Rune = '@'
	player.Speed = 1.0
	player.ActionPoints = 0
	player.SigntRange = 7

	levels := make(map[string]*Level)

	filenames, err := filepath.Glob("game/maps/*.map")
	if err != nil {
		panic(err)
	}
	for _, filename := range filenames {
		extIndex := strings.LastIndex(filename, ".map")
		lastSlashIndex := int(math.Max(float64(strings.LastIndex(filename, "\\")), float64(strings.LastIndex(filename, "/"))))
		levelName := filename[lastSlashIndex+1 : extIndex]
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
		level.Player = player
		level.Map = make([][]Tile, len(levelLines))
		level.Monsters = make(map[Pos]*Monster)
		level.Items = make(map[Pos][]*Item)
		level.Portals = make(map[Pos]*LevelPos)

		for i := range level.Map {
			level.Map[i] = make([]Tile, longestRow)
		}

		for y := 0; y < len(level.Map); y++ {
			line := levelLines[y]
			for x, c := range line {
				pos := Pos{x, y}
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
				case 's':
					level.Items[pos] = append(level.Items[pos], NewSword(pos))
					t.Rune = Pending
				case 'h':
					level.Items[pos] = append(level.Items[pos], NewHelmet(pos))
					t.Rune = Pending
				case '.':
					t.Rune = DirtFloor
				case '@':
					level.Player.X = x
					level.Player.Y = y
					t.Rune = Pending
				case 'R':
					level.Monsters[pos] = NewRat(pos)
					t.Rune = Pending
				case 'S':
					level.Monsters[pos] = NewSpider(pos)
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
				}
			}
		}

		levels[levelName] = level
	}

	return levels
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
		level.LastEvent = DoorOpen
		level.lineOfSight()
	}
}

func (game *Game) Move(to Pos) {
	level := game.CurrentLevel
	player := level.Player
	levelAndPos := level.Portals[to]
	if levelAndPos != nil {
		game.CurrentLevel = levelAndPos.Level
		game.CurrentLevel.Player.Pos = levelAndPos.Pos
		game.CurrentLevel.lineOfSight()
	} else {
		player.Pos = to

		/*
			items := level.Items[player.Pos]
			if len(items) > 0 {
				level.MoveItem(items[0], &player.Character)
				fmt.Println("Player inventory")
				for _, item := range player.Items {
					fmt.Println("			 ", item)
				}
			}

		*/

		level.LastEvent = Move
		for y, row := range level.Map {
			for x, _ := range row {
				level.Map[y][x].Visible = false
			}
		}
		level.lineOfSight()
	}
}

func (game *Game) resolveMovement(pos Pos) {
	level := game.CurrentLevel
	monster, exists := level.Monsters[pos]
	if exists {
		level.Attack(&level.Player.Character, &monster.Character)
		level.LastEvent = Attack
		if monster.Hitpoints <= 0 {
			monster.Kill(level)
		}
		if level.Player.Hitpoints <= 0 {
			panic("ded")
		}
	} else if canWalk(level, pos) {
		game.Move(pos)
	} else {
		checkDoor(level, pos)
	}
}

func (game *Game) handleInput(input *Input) {
	level := game.CurrentLevel
	p := level.Player

	switch input.Typ {
	case Up:
		newPos := Pos{p.X, p.Y - 1}
		game.resolveMovement(newPos)
	case Down:
		newPos := Pos{p.X, p.Y + 1}
		game.resolveMovement(newPos)
	case Left:
		newPos := Pos{p.X - 1, p.Y}
		game.resolveMovement(newPos)
	case Right:
		newPos := Pos{p.X + 1, p.Y}
		game.resolveMovement(newPos)
	case TakeAll:
		for _, item := range level.Items[p.Pos] {
			level.MoveItem(item, &level.Player.Character)
		}
		level.LastEvent = PickUp
	case TakeItem:
		level.MoveItem(input.Item, &level.Player.Character)
		level.LastEvent = PickUp
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

func (game *Game) Run() {
	fmt.Println("Starting game...")
	//Level := LoadLevelFromFile("game/maps/level1.map")
	for _, lchan := range game.LevelChans {
		lchan <- game.CurrentLevel
	}

	for input := range game.InputChan {
		if input != nil && input.Typ == QuitGame {
			fmt.Println("QUIT")
			return
		}
		//p := game.Level.Player.Pos
		//level.bresenham(p, Pos{p.X + 5, p.Y + 5})

		game.handleInput(input)
		for _, monster := range game.CurrentLevel.Monsters {
			monster.Update(game.CurrentLevel)
		}

		/*
			if len(game.LevelChans) == 0 {
				return
			}
		*/

		for _, lchan := range game.LevelChans {
			lchan <- game.CurrentLevel
		}
	}
}
