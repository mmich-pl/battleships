package app

import (
	"battleships/internal/models"
	. "battleships/internal/utils"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/sync/errgroup"
	"strconv"
	"strings"
	"time"
)

var (
	playerBoardPos   = [2]int{2, 5}
	opponentBoardPos = [2]int{50, 5}
	playerNickPos    = [2]int{2, 27}
	playerDescPos    = [2]int{2, 28}
	opponentNickPos  = [2]int{50, 27}
	opponentDescPos  = [2]int{50, 28}

	playerMovePos = [2]int{2, 3}
	turnPos       = [2]int{65, 3}
	timerPos      = [2]int{50, 3}
	gameResultPos = [2]int{2, 32}
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

func (a *App) markPlayerMove(opponentBoard *gui.Board, state gui.State, coord string) error {
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

		switch state := &a.PlayerBoardState[x][y]; *state {
		case gui.Hit, gui.Ship:
			*state = gui.Hit
		default:
			*state = gui.Miss
		}
	}
	playerBoard.SetStates(a.PlayerBoardState)
	return nil
}

func (a *App) drawBoard(ui gui.GUI, playerBoard, opponentBoard *gui.Board) {
	playerBoard.SetStates(a.PlayerBoardState)
	opponentBoard.SetStates(a.OpponentBoardState)

	a.renderDescription(ui, a.Description.Desc, a.Description.OpponentDescription)

	ui.Draw(playerBoard)
	ui.Draw(opponentBoard)
}

func (a *App) renderDescription(g gui.GUI, playerDescription, opponentDescription string) {
	g.Draw(gui.NewText(playerNickPos[0], playerNickPos[1], a.Description.Nick, nil))
	g.Draw(gui.NewText(opponentNickPos[0], opponentNickPos[1], a.Description.Opponent, nil))

	fragments := [2]struct {
		desc []string
		pos  [2]int
	}{
		{strings.Split(wordwrap.WrapString(playerDescription, 40), "\n"), playerDescPos},
		{strings.Split(wordwrap.WrapString(opponentDescription, 40), "\n"), opponentDescPos},
	}
	for _, frag := range fragments {
		for i, f := range frag.desc {
			g.Draw(gui.NewText(frag.pos[0], frag.pos[1]+i, f, &gui.TextConfig{
				FgColor: gui.White,
				BgColor: gui.Grey,
			}))
		}
	}
}

func (a *App) RenderBoards(status *models.StatusResponse) {
	ui := gui.NewGUI(true)
	playerBoard := gui.NewBoard(playerBoardPos[0], playerBoardPos[1], nil)
	opponentBoard := gui.NewBoard(opponentBoardPos[0], opponentBoardPos[1], nil)

	a.drawBoard(*ui, playerBoard, opponentBoard)
	playerMove := gui.NewText(playerMovePos[0], playerMovePos[1], "Press on any coordinate to take a shot.", nil)
	ui.Draw(playerMove)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(status.Timer)*time.Second)
	defer cancel()

	playerTurn := gui.NewText(turnPos[0], turnPos[1], fmt.Sprintf("Should I fire: %t", status.ShouldFire), nil)
	timer := gui.NewText(timerPos[0], timerPos[1], fmt.Sprintf("Timer: %d", status.Timer), nil)
	gameResult := gui.NewText(gameResultPos[0], gameResultPos[1], fmt.Sprintf("Game is running!"), nil)
	ui.Draw(playerTurn)
	ui.Draw(timer)
	ui.Draw(gameResult)

	//update timer
	go func() {
		for {
			status, _ = a.client.GameStatus(GameStatusEndpoint)
			time.Sleep(time.Second / 4)
			timer.SetText(fmt.Sprintf("Timer: %d", status.Timer))
			playerTurn.SetText(fmt.Sprintf(If(status.ShouldFire,
				"It's your turn, fire at will", "It's your opponent turn, be patient")))
			_ = a.markOpponentMoves(playerBoard, status)
		}
	}()

	//get input coords
	go func() {
		for {
			for status.ShouldFire == true {
				char := opponentBoard.Listen(ctx)
				shoot, _ := a.client.Fire(FireEndpoint, char)
				var state = If(shoot.Result == "hit" || shoot.Result == "sunk", gui.Hit, gui.Miss)
				_ = a.markPlayerMove(opponentBoard, state, char)
				playerMove.SetText(fmt.Sprintf("Last player move: %s", char))
			}
		}
	}()

	// handle end game
	go func() {
		for {
			if status.GameStatus == "ended" {
				gameResult.SetText(If(status.LastGameStatus == "win",
					"Game ended, You win", "Game ended, You lost"))
			}
		}
	}()

	ui.Start(nil)
}
