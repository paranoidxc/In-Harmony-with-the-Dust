package server

import (
	"context"
	"errors"
	"fmt"
	"multiplayer-game-with-go/pkg/backend"
	"multiplayer-game-with-go/proto"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	clientTimeout = 15
	maxClients    = 8
)

type client struct {
	streamServer proto.Game_StreamServer
	lastMessage  time.Time
	done         chan error
	playerID     uuid.UUID
	id           uuid.UUID
}

type GameServer struct {
	game     *backend.Game
	mu       sync.RWMutex
	password string
}

func NewGameServer(game *backend.Game, password string) *GameServer {
	server := &GameServer{
		game:     game,
		password: password,
	}
	return server
}

func (s *GameServer) Stream(srv proto.Game_StreamServer) error {
	fmt.Println("stream")
	return nil
}

func (s *GameServer) Connect(ctx context.Context, req *proto.ConnectRequest) (*proto.ConnectResponse, error) {
	fmt.Println("connect")

	playerID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	fmt.Println("playerID", playerID)

	// Exit as early as possible if password is wrong.
	if req.Password != s.password {
		return nil, errors.New("invalid password provided")
	}
	fmt.Println("req", req)

	s.game.Mu.RLock()
	if s.game.GetEntity(playerID) != nil {
		return nil, errors.New("duplicate player ID provided")
	}
	s.game.Mu.RUnlock()

	player := &backend.Player{}

	s.game.Mu.Lock()
	s.game.AddEntity(player)
	s.game.Mu.Unlock()

	// Inform all other clients of the new player.
	//s.game.Mu.RLock()
	//entities := make([]*proto.Entity, 0)

	// Add the new client.
	s.mu.Lock()
	token := uuid.New()
	s.mu.Unlock()

	return &proto.ConnectResponse{
		Token: token.String(),
	}, nil
}
