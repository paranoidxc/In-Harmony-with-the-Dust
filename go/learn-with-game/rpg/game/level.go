package game

import (
	"math"
	"strconv"
	"time"
)

type Level struct {
	Map       [][]Tile
	Player    *Player
	Monsters  map[Pos]*Monster
	Items     map[Pos][]*Item
	Portals   map[Pos]*LevelPos
	Events    []string
	Debug     map[Pos]bool
	EventPos  int
	LastEvent GameEvent
}

func (level *Level) MoveItem(itemToMove *Item, character *Character) {
	pos := character.Pos
	items := level.Items[pos]
	for i, item := range items {
		if item == itemToMove {
			items = append(items[:i], items[i+1:]...)
			level.Items[pos] = items
			character.Items = append(character.Items, item)
			level.AddEvent(character.Name + " pick up:" + item.Name)
			return
		}
	}
}

func (level *Level) AddEvent(event string) {
	level.Events[level.EventPos] = event

	level.EventPos++
	if level.EventPos == len(level.Events) {
		level.EventPos = 0
	}
}

func (level *Level) Attack(c1 *Character, c2 *Character) {
	c1.ActionPoints--
	c1AttackPower := c1.Strength
	c2.Hitpoints -= c1AttackPower

	if c2.Hitpoints > 0 {
		level.AddEvent(c1.Name + " Attacked " + c2.Name + " for " + strconv.Itoa(c1AttackPower))
	}
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
