package client

import (
	"context"
	"fmt"
	"log"
	"multiplayer-game-with-go/pkg/backend"
	"multiplayer-game-with-go/pkg/frontend"
	"multiplayer-game-with-go/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const (
	positionHistoryLimit = 5
)

type GameClient struct {
	CurrentPlayer   uuid.UUID
	Stream          proto.Game_StreamClient
	Game            *backend.Game
	View            *frontend.View
	positionHistory []backend.Coordinate
}

func NewGameClient(game *backend.Game, view *frontend.View) *GameClient {
	return &GameClient{
		Game:            game,
		View:            view,
		positionHistory: make([]backend.Coordinate, positionHistoryLimit),
	}
}

func (c *GameClient) Connect(grpcClient proto.GameClient, playerID uuid.UUID, playerName string, password string) error {
	req := proto.ConnectRequest{
		Id:       playerID.String(),
		Name:     playerName,
		Password: password,
	}
	resp, err := grpcClient.Connect(context.Background(), &req)
	if err != nil {
		return err
	}
	log.Println("grpc resp", resp)

	// Add initial entity state.
	for _, entity := range resp.Entities {
		backendEntity := proto.GetBackendEntity(entity)
		if backendEntity == nil {
			return fmt.Errorf("can not get backend entity from %+v", entity)
		}
		c.Game.AddEntity(backendEntity)
	}
	// Initialize stream with token.
	header := metadata.New(map[string]string{"authorization": resp.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), header)
	stream, err := grpcClient.Stream(ctx)
	if err != nil {
		return err
	}

	c.CurrentPlayer = playerID
	c.View.CurrentPlayer = playerID
	c.Stream = stream

	return nil
}

// Exit stops the tview application and prints a message.
// This is needed as stdout is mangled while tview is running.
func (c *GameClient) Exit(message string) {
	c.View.App.Stop()
	log.Println("exit", message)
}

// Start begins the goroutines needed to recieve server changes and send game
// changes.
func (c *GameClient) Start() {
	// Handle local game engine changes.
	go func() {
		for {
			change := <-c.Game.ChangeChannel
			switch change.(type) {
			case backend.MoveChange:
				log.Print("change", change)
			}
		}
	}()

	// Handle stream messages.
	go func() {
		for {
			resp, err := c.Stream.Recv()
			if err != nil {
				c.Exit(fmt.Sprintf("can not receive, error: %v", err))
				return
			}

			log.Println(resp)
		}
	}()
}
