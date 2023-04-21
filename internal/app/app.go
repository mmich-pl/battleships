package app

import (
	"battleships/internal/battlehip_client"
	"fmt"
)

const (
	InitEndpoint = "/game"
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
	return nil
}
