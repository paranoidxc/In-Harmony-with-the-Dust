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
				log.Print("movechange", change)
				change := change.(backend.MoveChange)
				c.handleMoveChange(change)
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

			log.Println("client stream get recv", resp)

			c.Game.Mu.Lock()
			switch resp.GetAction().(type) {
			case *proto.Response_AddEntity:
				c.handleAddEntityResponse(resp)
			case *proto.Response_UpdateEntity:
				c.handleUpdateEntityResponse(resp)
				/*
					case *proto.Response_RemoveEntity:
						c.handleRemoveEntityResponse(resp)
					case *proto.Response_PlayerRespawn:
						c.handlePlayerRespawnResponse(resp)
					case *proto.Response_RoundOver:
						c.handleRoundOverResponse(resp)
					case *proto.Response_RoundStart:
						c.handleRoundStartResponse(resp)
				*/
			default:
				log.Panicln("TTTTTTTTTDDDDDDDDDD")
			}
			c.Game.Mu.Unlock()
		}
	}()
}

func (c *GameClient) handleMoveChange(change backend.MoveChange) {
	req := proto.Request{
		Action: &proto.Request_Move{
			Move: &proto.Move{
				Direction: proto.GetProtoDirection(change.Direction),
			},
		},
	}
	c.Stream.Send(&req)
	// Store position history to help with stuttering.
	c.positionHistory = append([]backend.Coordinate{change.Position}, c.positionHistory[:positionHistoryLimit]...)
}

func (c *GameClient) handleAddEntityResponse(resp *proto.Response) {
	add := resp.GetAddEntity()
	entity := proto.GetBackendEntity(add.Entity)
	if entity == nil {
		c.Exit(fmt.Sprintf("can not get backend entity from %+v", entity))
		return
	}
	// To prevent jittering, ignore lasers we created.
	/*
		laser, ok := entity.(*backend.Laser)
		if ok && laser.OwnerID == c.CurrentPlayer {
			return
		}
	*/
	c.Game.AddEntity(entity)
}
func (c *GameClient) handleUpdateEntityResponse(resp *proto.Response) {
	update := resp.GetUpdateEntity()
	entity := proto.GetBackendEntity(update.Entity)
	if entity == nil {
		c.Exit(fmt.Sprintf("can not get backend entity from %+v", entity))
		return
	}
	// To prevent jittering, ignore updates for recent positions.
	// Note: This feels OK, but isn't perfect. I think if I refactored the
	// networking to use more targeted responses, i.e. "move confirmed" sent
	// after a player moves, you could compare recent moves instead of
	// positions and do something like rollback networking.
	player, ok := entity.(*backend.Player)
	if ok && player.ID() == c.CurrentPlayer {
		for _, position := range c.positionHistory {
			if player.Position() == position {
				return
			}
		}
	}
	c.Game.UpdateEntity(entity)
}
