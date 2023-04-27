package app

import (
	"battleships/internal/battlehip_client"
	"battleships/internal/models"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"log"
	"time"
)

const (
	InitEndpoint        = "/game"
	GameStatusEndpoint  = "/game"
	BoardEndpoint       = "/game/board"
	OpponentDescription = "/game/desc"
)

type App struct {
	client             *battlehip_client.BattleshipHTTPClient
	PlayerBoardState   [10][10]gui.State
	OpponentBoardState [10][10]gui.State
}

func New(c *battlehip_client.BattleshipHTTPClient) *App {
	return &App{
		client: c,
	}
}

func (a *App) Run() error {
	err := a.client.InitGame(InitEndpoint, "", "", "", true)
	if err != nil {
		return fmt.Errorf("failed to init game: %w", err)
	}

	status, err := a.waitForGameStart(err)
	if err != nil {
		return fmt.Errorf("failed to get game status: %w", err)
	}

	desc, err := a.client.Descriptions(OpponentDescription)
	if err != nil {
		return fmt.Errorf("failed to get game status: %w", err)
	}
	log.Print(status)

	board, err := a.client.Board(BoardEndpoint)
	if err != nil {
		return fmt.Errorf("failed to get board: %w", err)
	}

	if p, o, err := a.setUpBoardsState(board); err != nil {
		return err
	} else {
		a.OpponentBoardState = *o
		a.PlayerBoardState = *p
	}

	RenderBoards(status, desc, a.PlayerBoardState, a.OpponentBoardState)
	return nil
}

func (a *App) waitForGameStart(err error) (*models.StatusResponse, error) {
	status, err := a.client.GameStatus(GameStatusEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get game status: %w", err)
	}

	for status.GameStatus != "game_in_progress" {
		time.Sleep(time.Second)
		status, err = a.client.GameStatus(GameStatusEndpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to get game status: %w", err)
		}
	}
	return status, nil
}
