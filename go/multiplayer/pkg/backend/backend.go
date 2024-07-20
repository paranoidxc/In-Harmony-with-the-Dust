package backend

import (
	"fmt"
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

// Identifier is an entity that provides an ID method.
type Identifier interface {
	ID() uuid.UUID
}

type Game struct {
	Entities      map[uuid.UUID]Identifier
	Mu            sync.RWMutex
	ActionChannel chan Action
}

func NewGame() *Game {
	return &Game{
		Entities:      make(map[uuid.UUID]Identifier),
		ActionChannel: make(chan Action, 1),
	}
}

func (game *Game) Start() {
	go game.watchActions()
	go game.watchCollisions()
}

func (game *Game) watchActions() {
	for {
		action := <-game.ActionChannel
		fmt.Println("watchCollisions", action)
	}
}

func (game *Game) watchCollisions() {
	for {
		//fmt.Println("watchCollisions", collisionCheckFrequency)
		time.Sleep(collisionCheckFrequency)
	}
}

func (game *Game) GetEntity(id uuid.UUID) Identifier {
	return game.Entities[id]
}

func (game *Game) AddEntity(entity Identifier) {
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
