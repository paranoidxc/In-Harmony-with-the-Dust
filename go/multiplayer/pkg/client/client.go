package client

import (
	"context"
	"fmt"
	"multiplayer-game-with-go/pkg/backend"
	"multiplayer-game-with-go/pkg/frontend"
	"multiplayer-game-with-go/proto"

	"github.com/google/uuid"
)

type GameClient struct {
}

func NewGameClient(game *backend.Game, view *frontend.View) *GameClient {
	return &GameClient{}
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
	fmt.Println("grpc resp", resp)

	return nil
}
