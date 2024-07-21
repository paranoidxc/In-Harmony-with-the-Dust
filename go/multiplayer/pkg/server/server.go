package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"multiplayer-game-with-go/pkg/backend"
	"multiplayer-game-with-go/proto"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc/metadata"
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
	clients  map[uuid.UUID]*client
	mu       sync.RWMutex
	password string
}

func NewGameServer(game *backend.Game, password string) *GameServer {
	server := &GameServer{
		game:     game,
		clients:  make(map[uuid.UUID]*client),
		password: password,
	}

	server.watchChanges()
	return server
}

func (s *GameServer) getClientFromContext(ctx context.Context) (*client, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	tokenRaw := headers["authorization"]
	if len(tokenRaw) == 0 {
		return nil, errors.New("no token provided")
	}
	token, err := uuid.Parse(tokenRaw[0])
	if err != nil {
		return nil, errors.New("cannot parse token")
	}
	s.mu.RLock()
	currentClient, ok := s.clients[token]
	s.mu.RUnlock()
	if !ok {
		return nil, errors.New("token not recognized")
	}
	return currentClient, nil
}

func (s *GameServer) Stream(srv proto.Game_StreamServer) error {
	//fmt.Println("stream")
	ctx := srv.Context()
	currentClient, err := s.getClientFromContext(ctx)
	if err != nil {
		return err
	}
	if currentClient.streamServer != nil {
		return errors.New("stream already active")
	}
	currentClient.streamServer = srv

	log.Println("start new server")

	// Wait for stream requests.
	go func() {
		for {
			req, err := srv.Recv()
			if err != nil {
				log.Printf("receive error %v", err)
				currentClient.done <- errors.New("failed to receive request")
				return
			}
			log.Printf("got message\n")
			log.Println(req)
			currentClient.lastMessage = time.Now()

			switch req.GetAction().(type) {
			case *proto.Request_Move:
				log.Println("Request_Move")
				s.handleMoveRequest(req, currentClient)
			case *proto.Request_Laser:
				log.Println("Request_Laser")
				//s.handleLaserRequest(req, currentClient)
			}
		}
	}()

	var doneError error
	select {
	case <-ctx.Done():
		doneError = ctx.Err()
	case doneError = <-currentClient.done:
		fmt.Println("currentClient.done")
	}

	log.Printf(`stream done with error "%v"`, doneError)
	log.Printf("%s - removing client", currentClient.id)

	return nil
}

func (s *GameServer) Connect(ctx context.Context, req *proto.ConnectRequest) (*proto.ConnectResponse, error) {
	log.Println("connect")

	playerID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	fmt.Println("playerID", playerID)

	// Exit as early as possible if password is wrong.
	if req.Password != s.password {
		return nil, errors.New("invalid password provided")
	}
	log.Println("req", req)

	s.game.Mu.RLock()
	if s.game.GetEntity(playerID) != nil {
		return nil, errors.New("duplicate player ID provided")
	}
	s.game.Mu.RUnlock()

	re := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !re.MatchString(req.Name) {
		return nil, errors.New("invalid name provided")
	}
	icon, _ := utf8.DecodeRuneInString(strings.ToUpper(req.Name))

	// Choose a random spawn point.
	spawnPoints := s.game.GetMapByType()[backend.MapTypeSpawn]
	i := rand.Int() % len(spawnPoints)
	startCoordinate := spawnPoints[i]

	player := &backend.Player{
		Name:            req.Name,
		Icon:            icon,
		IdentifierBase:  backend.IdentifierBase{UUID: playerID},
		CurrentPosition: startCoordinate,
	}

	s.game.Mu.Lock()
	s.game.AddEntity(player)
	s.game.Mu.Unlock()

	// Inform all other clients of the new player.
	s.game.Mu.RLock()
	entities := make([]*proto.Entity, 0)
	for _, entity := range s.game.Entities {
		protoEntity := proto.GetProtoEntity(entity)
		if protoEntity != nil {
			entities = append(entities, protoEntity)
		}
	}
	s.game.Mu.RUnlock()

	// Inform all other clients of the new player.
	resp := proto.Response{
		Action: &proto.Response_AddEntity{
			AddEntity: &proto.AddEntity{
				Entity: proto.GetProtoEntity(player),
			},
		},
	}
	s.broadcast(&resp)

	// Add the new client.
	s.mu.Lock()
	token := uuid.New()
	s.clients[token] = &client{
		id:          token,
		playerID:    playerID,
		done:        make(chan error),
		lastMessage: time.Now(),
	}
	s.mu.Unlock()

	return &proto.ConnectResponse{
		Token:    token.String(),
		Entities: entities,
	}, nil
}

func (s *GameServer) watchChanges() {
	go func() {
		for {
			change := <-s.game.ChangeChannel
			log.Printf("ChangeChannel %+v", change)
			switch change.(type) {
			case backend.MoveChange:
				change := change.(backend.MoveChange)
				s.handleMoveChange(change)
			default:
				panic("FFFFFFF")
			}
		}
	}()
}

// handleMoveRequest makes a request to the game engine to move a player.
func (s *GameServer) handleMoveRequest(req *proto.Request, currentClient *client) {
	move := req.GetMove()
	s.game.ActionChannel <- backend.MoveAction{
		ID:        currentClient.playerID,
		Direction: proto.GetBackendDirection(move.Direction),
		Created:   time.Now(),
	}
}

func (s *GameServer) handleMoveChange(change backend.MoveChange) {
	resp := proto.Response{
		Action: &proto.Response_UpdateEntity{
			UpdateEntity: &proto.UpdateEntity{
				Entity: proto.GetProtoEntity(change.Entity),
			},
		},
	}
	s.broadcast(&resp)
}

// broadcast sends a response to all clients.
func (s *GameServer) broadcast(resp *proto.Response) {
	s.mu.Lock()
	for id, currentClient := range s.clients {
		if currentClient.streamServer == nil {
			continue
		}
		if err := currentClient.streamServer.Send(resp); err != nil {
			log.Printf("%s - broadcast error %v", id, err)
			currentClient.done <- errors.New("failed to broadcast message")
			continue
		}
		log.Printf("%s - broadcasted %+v", resp, id)
	}
	s.mu.Unlock()
}
