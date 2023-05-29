package app

import (
	"battleships/internal/app/menu"
	"battleships/internal/battleship_client"
	"battleships/internal/models"
	. "battleships/internal/utils"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"strings"
	"time"
)

type App struct {
	Client             battleship_client.BattleshipClient
	PlayerBoardState   [10][10]gui.State
	OpponentBoardState [10][10]gui.State
	Description        *models.DescriptionResponse
}

func New(c battleship_client.BattleshipClient) *App {
	return &App{
		Client: c,
	}
}

func (a *App) Run() error {
	if a.Client.GetToken() != "" {
		a.Client.ResetToken()
	}

	var playerNick string
	var playerDescription string

	if a.Description != nil && a.Description.Nick != "" {
		playerNick = a.Description.Nick
		playerDescription = a.Description.Desc
	} else {
		nickIsValid := false
		for !nickIsValid {
			playerNick, _ = GetPlayerInput("set your nickname or hit enter to get autogenerated one: ")
			nickLen := len(playerNick)
			if nickLen == 0 {
				nickIsValid = true
			} else if nickLen < 2 || nickLen > 20 {
				nickIsValid = false
			} else {
				nickIsValid = true
			}
		}
		playerDescription, _ = GetPlayerInput("set your description or hit enter to get autogenerated one: ")
	}

	fleet, _ := GetPlayerInput("do you want to set your own fleet? [y/n]: ")
	ownFleet := If(strings.ToLower(fleet) == "y", true, false)
	var playerShips []string
	if ownFleet {
		playerShips = RenderInputBoard()
	}
	playerShipsCoordsList := If(len(playerShips) != 0, playerShips, nil)

	bot, _ := GetPlayerInput("do you want to play with bot? [y/n]: ")
	playWithBot := If(strings.ToLower(bot) == "y", true, false)
	enemyNick := ""
	if !playWithBot {
		fmt.Println("Type player name from list:")
		err := menu.ListPlayer(a.Client)
		fmt.Println()

		if err != nil {
			return err
		}
		enemyNick, _ = GetPlayerInput("type enemy nick or leave blank to wait for other players request: ")
	}
	err := a.Client.InitGame(playerNick, playerDescription, enemyNick, playerShipsCoordsList, playWithBot)
	if err != nil {
		return fmt.Errorf("failed to init game: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println("the game has been initiated, waiting for opponent")
	go func(ctx context.Context) {
		for {
			time.Sleep(10 * time.Second)
			select {
			case <-ctx.Done():
				return
			default:
				_ = a.Client.RefreshSession()
			}
		}
	}(ctx)

	status, err := a.setUpGame(cancel)
	if err != nil {
		return err
	}

	bd := InitBoardData(a)

	err = bd.RenderGameBoards(status)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) setUpGame(cancel context.CancelFunc) (*models.StatusResponse, error) {
	status, err := a.waitForGameStart(cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to get game status: %w", err)
	}

	a.Description, err = a.Client.Description()
	if err != nil {
		return nil, fmt.Errorf("failed to get game status: %w", err)
	}

	board, err := a.Client.Board()
	if err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	if err = a.setUpBoardsState(board); err != nil {
		return nil, err
	}

	return status, nil
}

func (a *App) waitForGameStart(cancel context.CancelFunc) (*models.StatusResponse, error) {
	status, err := a.Client.GameStatus()

	if err != nil {
		return nil, fmt.Errorf("failed to get game status: %w", err)
	}

	type channelResponse struct {
		*models.StatusResponse
		error
	}

	ticker := time.NewTicker(time.Second)
	channel := make(chan channelResponse, 1)

	go func() {
		for range ticker.C {
			if status.GameStatus == "game_in_progress" {
				channel <- channelResponse{status, nil}
				cancel()
				break
			}

			status, err = a.Client.GameStatus()

			if err != nil {
				cancel()
				channel <- channelResponse{nil, err}
			}
		}
	}()

	resp := <-channel
	return resp.StatusResponse, resp.error
}
