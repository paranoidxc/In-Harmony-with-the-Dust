package backend

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	roundOverScore          = 10
	newRoundWaitTime        = 10 * time.Second
	collisionCheckFrequency = 10 * time.Millisecond
	moveThrottle            = 100 * time.Millisecond
	laserThrottle           = 500 * time.Millisecond
	laserSpeed              = 50
)

type Action interface {
	Perform(game *Game)
}

type MoveAction struct {
	Direction Direction
	ID        uuid.UUID
	Created   time.Time
}

// Perform contains backend logic required to move an entity.
func (action MoveAction) Perform(game *Game) {
	log.Println("moveaction", action)
	entity := game.GetEntity(action.ID)
	if entity == nil {
		return
	}

	mover, ok := entity.(Mover)
	if !ok {
		return
	}

	positioner, ok := entity.(Positioner)
	if !ok {
		return
	}

	//actionKey := fmt.Sprintf("%T:%s", action, entity.ID().String())

	position := positioner.Position()
	// Move the entity.
	switch action.Direction {
	case DirectionUp:
		position.Y--
	case DirectionDown:
		position.Y++
	case DirectionLeft:
		position.X--
	case DirectionRight:
		position.X++
	}

	mover.Move(position)

	// Inform the client that the entity moved.
	change := MoveChange{
		Entity:    entity,
		Direction: action.Direction,
		Position:  position,
	}
	game.sendChange(change)
	//game.updateLastActionTime(actionKey, action.Created)
}

type Change interface{}

// MoveChange is sent when the game engine moves an entity.
type MoveChange struct {
	Change
	Entity    Identifier
	Direction Direction
	Position  Coordinate
}

// Identifier is an entity that provides an ID method.
type Identifier interface {
	ID() uuid.UUID
}

type Game struct {
	Entities      map[uuid.UUID]Identifier
	gameMap       [][]rune
	Mu            sync.RWMutex
	ActionChannel chan Action
	ChangeChannel chan Change
	WaitForRound  bool
}

func NewGame() *Game {
	return &Game{
		Entities:      make(map[uuid.UUID]Identifier),
		ActionChannel: make(chan Action, 1),
		ChangeChannel: make(chan Change, 1),
		gameMap:       MapDefault,
		WaitForRound:  false,
	}
}

func (game *Game) Start() {
	go game.watchActions()
	go game.watchCollisions()
}

func (game *Game) watchActions() {
	for {
		action := <-game.ActionChannel
		log.Println("watchActions", action)

		if game.WaitForRound {
			continue
		}

		game.Mu.Lock()
		action.Perform(game)
		game.Mu.Unlock()
	}
}

func (game *Game) watchCollisions() {
	for {
		//fmt.Println("watchCollisions", collisionCheckFrequency)
		time.Sleep(collisionCheckFrequency)
	}
}

// sendChange sends a change to the change channel.
func (game *Game) sendChange(change Change) {
	select {
	case game.ChangeChannel <- change:
	default:
	}
}

func (game *Game) GetEntity(id uuid.UUID) Identifier {
	return game.Entities[id]
}

func (game *Game) AddEntity(entity Identifier) {
	game.Entities[entity.ID()] = entity
}

// UpdateEntity updates an entity.
func (game *Game) UpdateEntity(entity Identifier) {
	game.Entities[entity.ID()] = entity
}

// IdentifierBase is embedded to satisfy the Identifier interface.
type IdentifierBase struct {
	UUID uuid.UUID
}

// ID returns the UUID of an entity.
func (e IdentifierBase) ID() uuid.UUID {
	return e.UUID
}

// Coordinate is used for all position-related variables.
type Coordinate struct {
	X int
	Y int
}

// Direction is used to represent Direction constants.
type Direction int

// Contains direction constants - DirectionStop will take no effect.
const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
	DirectionStop
)

// Positioner is an entity that has a position.
type Positioner interface {
	Position() Coordinate
}

// Mover is an entity that can be moved.
type Mover interface {
	Move(Coordinate)
}
