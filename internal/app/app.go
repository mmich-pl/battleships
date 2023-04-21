package app

import (
	"battleships/internal/battlehip_client"
	"fmt"
	"log"
)

const (
	InitEndpoint       = "/game"
	GameStatusEndpoint = "/game"
)

type App struct {
	client *battlehip_client.BattleshipHTTPClient
}

func New(c *battlehip_client.BattleshipHTTPClient) *App {
	return &App{
		client: c,
	}
}
func (a *App) Run() error {
	err := a.client.InitGame(InitEndpoint, "", "", "", false)
	if err != nil {
		return fmt.Errorf("failed to init game: %w", err)
	}

	status, err := a.client.GameStatus(GameStatusEndpoint)
	if err != nil {
		return fmt.Errorf("failed to get game status: %w", err)
	}

	log.Print(status)
	return nil
}
