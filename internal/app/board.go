package app

import (
	"battleships/internal/models"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/sync/errgroup"
	"log"
	"strconv"
	"strings"
	"time"
)

// Parses coordinates to two integers that represents board square in matrix
func mapCoords(coordinate string) (int, int, error) {
	column := coordinate[0]
	if 'A' > column || column > 'J' {
		return -1, -1, fmt.Errorf("coordinate out of bound: expected column in bounds [A-J]")
	}

	x := int(column - 'A')
	y, err := strconv.Atoi(coordinate[1:])
	if err != nil {
		return -1, -1, fmt.Errorf("wrong coordinates format: %w", err)
	} else if y < 1 || y > 10 {
		return -1, -1, fmt.Errorf("coordinate [%s] out of bound: expected row in bounds [1-10]", coordinate)
	}
	return x, y - 1, nil
}

// Parses API response to [10][10] matrix format used by client
func (a *App) setUpBoardsState(board []string) error {
	for i := 0; i < len(a.PlayerBoardState); i++ {
		a.PlayerBoardState[i] = [10]gui.State{}
		a.OpponentBoardState[i] = [10]gui.State{}
	}

	g := new(errgroup.Group)
	for _, coords := range board {
		c := coords
		g.Go(func() error {
			if x, y, err := mapCoords(c); err != nil {
				return err
			} else {
				a.PlayerBoardState[x][y] = gui.Ship
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to fill boards: %w", err)
	}

	return nil
}

func renderDescription(g gui.GUI, playerDescription, opponentDescription string) {
	fragments := [2]struct {
		desc []string
		x    int
		y    int
	}{
		{strings.Split(wordwrap.WrapString(playerDescription, 40), "\n"), 2, 27},
		{strings.Split(wordwrap.WrapString(opponentDescription, 40), "\n"), 50, 27},
	}
	for _, frag := range fragments {
		for i, f := range frag.desc {
			g.Draw(gui.NewText(frag.x, frag.y+i, f, &gui.TextConfig{
				FgColor: gui.White,
				BgColor: gui.Grey,
			}))
		}
	}
}

func (a *App) drawBoard(ui gui.GUI, playerBoard, opponentBoard *gui.Board) {
	playerBoard.SetStates(a.PlayerBoardState)
	opponentBoard.SetStates(a.OpponentBoardState)

	ui.Draw(gui.NewText(2, 1, fmt.Sprintf("%s vs %s", a.Description.Nick, a.Description.Opponent), nil))
	renderDescription(ui, a.Description.Desc, a.Description.OpponentDescription)

	ui.Draw(playerBoard)
	ui.Draw(opponentBoard)
}

func (a *App) markHitOrMiss(opponentBoard *gui.Board, state gui.State, coord string) error {
	x, y, err := mapCoords(coord)
	if err != nil {
		return fmt.Errorf("failed to parse coord: %w", err)
	}

	a.OpponentBoardState[x][y] = state
	opponentBoard.SetStates(a.OpponentBoardState)
	return nil
}

func (a *App) markOpponentMoves(playerBoard *gui.Board, status *models.StatusResponse) error {

	for _, cords := range status.OpponentShots {
		x, y, err := mapCoords(cords)
		if err != nil {
			return fmt.Errorf("failed to parse coords: %w", err)
		}

		if a.PlayerBoardState[x][y] == gui.Ship || a.PlayerBoardState[x][y] == gui.Hit {
			a.PlayerBoardState[x][y] = gui.Hit
		} else {
			a.PlayerBoardState[x][y] = gui.Miss
		}

	}
	playerBoard.SetStates(a.PlayerBoardState)

	return nil
}

func (a *App) RenderBoards(status *models.StatusResponse) {
	ui := gui.NewGUI(true)
	playerBoard := gui.NewBoard(2, 5, nil)
	opponentBoard := gui.NewBoard(50, 5, nil)

	a.drawBoard(*ui, playerBoard, opponentBoard)
	playerMove := gui.NewText(2, 4, "Press on any coordinate to log it.", nil)
	ui.Draw(playerMove)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(status.Timer)*time.Second)
	//defer cancel()

	playerTurn := gui.NewText(2, 2, fmt.Sprintf("Should I fire: %t", status.ShouldFire), nil)
	timer := gui.NewText(2, 3, fmt.Sprintf("Timer: %d", status.Timer), nil)
	gameResult := gui.NewText(2, 32, fmt.Sprintf("Game has started!"), nil)
	ui.Draw(playerTurn)
	ui.Draw(timer)
	ui.Draw(gameResult)

	//update timer
	go func() {
		for {
			status, _ = a.client.GameStatus(GameStatusEndpoint)
			time.Sleep(time.Second / 4)
			timer.SetText(fmt.Sprintf("Timer: %d", status.Timer))
			playerTurn.SetText(fmt.Sprintf("Should I fire: %t", status.ShouldFire))
			_ = a.markOpponentMoves(playerBoard, status)
		}
	}()

	//get input coords
	go func() {
		for {
			for status.ShouldFire == true {
				char := opponentBoard.Listen(context.TODO())
				shoot, err := a.client.Fire(FireEndpoint, char)

				if err != nil {
					log.Fatalf("failed to perform fire request:%s", err)
				}

				var state gui.State
				if shoot.Result == "hit" || shoot.Result == "sunk" {
					state = gui.Hit
				} else if shoot.Result == "miss" {
					state = gui.Miss
				}
				_ = a.markHitOrMiss(opponentBoard, state, char)
				playerMove.SetText(fmt.Sprintf("Fired at: %s", char))
			}
		}
	}()

	// handle end game
	go func() {
		for {
			if status.GameStatus == "ended" {
				if status.LastGameStatus == "win" {
					gameResult.SetText("Game ended, You win")
				} else {
					gameResult.SetText("Game ended, You lost")
				}
			}
		}
	}()

	ui.Start(nil)
}
