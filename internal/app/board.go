package app

import (
	"battleships/internal/models"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/sync/errgroup"
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
	ui.Draw(playerTurn)
	ui.Draw(timer)

	go func() {
		for {
			status, _ = a.client.GameStatus(GameStatusEndpoint)
			time.Sleep(time.Second)
			timer.SetText(fmt.Sprintf("Timer: %d", int(status.Timer)))
			playerTurn.SetText(fmt.Sprintf("Should I fire: %t", status.ShouldFire))
			//gA.statusBoard.SetText(status.GameStatus)
			//_ = a.MarkOpponentMoves(playerBoard, status)
		}
	}()

	/*	go func() {
		for {
			//updatePlayerBoard(*playerBoard, &playerState, status.OpponentShots)
			//ui.Draw(playerBoard)
			if status.ShouldFire {
				//log.Println(status.ShouldFire)
				char := opponentBoard.Listen(ctx)
				playerMove.SetText(fmt.Sprintf("Ready! Aim at %s! FIRE!", char))
				ui.Log("Coordinate: %s", char) // logs are displayed after the game exits
			}

		}
	}()*/

	ui.Start(nil)
}
