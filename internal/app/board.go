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
	playerDescPos   = [2]int{2, 28}
	opponentDescPos = [2]int{50, 28}

	errorLog = gui.NewText(2, 35, "", nil)
)

type BoardData struct {
	app             *App
	ui              *gui.GUI
	playerBoard     *gui.Board
	opponentBoard   *gui.Board
	playerNick      *gui.Text
	opponentNick    *gui.Text
	playerTurn      *gui.Text
	timer           *gui.Text
	gameResult      *gui.Text
	statusAfterFire *gui.Text
	acc             *gui.Text
	playerMove      *gui.Text
}

func InitBoardData(a *App) *BoardData {
	return &BoardData{
		app:             a,
		ui:              gui.NewGUI(true),
		playerBoard:     gui.NewBoard(2, 5, nil),
		opponentBoard:   gui.NewBoard(50, 5, nil),
		playerNick:      gui.NewText(2, 27, a.Description.Nick, nil),
		opponentNick:    gui.NewText(50, 27, a.Description.Opponent, nil),
		playerTurn:      gui.NewText(65, 3, "", nil),
		timer:           gui.NewText(50, 3, "", nil),
		gameResult:      gui.NewText(2, 32, fmt.Sprintf("Game is running!"), nil),
		statusAfterFire: gui.NewText(2, 31, "", nil),
		acc:             gui.NewText(2, 33, "", nil),
		playerMove:      gui.NewText(2, 3, "Press on any coordinate to take a shot.", nil),
	}
}

// Parses coordinates to two integers that represents board square in matrix
func mapCoords(coordinate string) (int, int, error) {
	if len(coordinate) == 0 {
		return -1, -1, fmt.Errorf("coordinate is empty")
	}
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

func (bd *BoardData) markPlayerMove(state gui.State, coord string) error {
	x, y, err := mapCoords(coord)
	if err != nil {
		return fmt.Errorf("failed to parse coord: %w", err)
	}

	bd.app.OpponentBoardState[x][y] = state
	bd.opponentBoard.SetStates(bd.app.OpponentBoardState)
	return nil
}

func (bd *BoardData) markOpponentMoves(status *models.StatusResponse) error {

	for _, cords := range status.OpponentShots {
		x, y, err := mapCoords(cords)
		if err != nil {
			return fmt.Errorf("failed to parse coords: %w", err)
		}

		switch state := &bd.app.PlayerBoardState[x][y]; *state {
		case gui.Hit, gui.Ship:
			*state = gui.Hit
		default:
			*state = gui.Miss
		}
	}
	bd.playerBoard.SetStates(bd.app.PlayerBoardState)
	return nil
}

func (bd *BoardData) drawBoard() {
	bd.playerBoard.SetStates(bd.app.PlayerBoardState)
	bd.opponentBoard.SetStates(bd.app.OpponentBoardState)

	bd.ui.Draw(bd.playerNick)
	bd.ui.Draw(bd.opponentNick)

	bd.renderDescription()

	bd.ui.Draw(bd.playerBoard)
	bd.ui.Draw(bd.opponentBoard)

	bd.ui.Draw(bd.playerTurn)
	bd.ui.Draw(bd.timer)
	bd.ui.Draw(bd.gameResult)
	bd.ui.Draw(bd.statusAfterFire)
	bd.ui.Draw(bd.playerMove)
	bd.ui.Draw(errorLog)
}

func (bd *BoardData) renderDescription() {

	fragments := [2]struct {
		desc []string
		pos  [2]int
	}{
		{strings.Split(wordwrap.WrapString(bd.app.Description.Desc, 40), "\n"), playerDescPos},
		{strings.Split(wordwrap.WrapString(bd.app.Description.OpponentDescription, 40), "\n"), opponentDescPos},
	}
	for _, frag := range fragments {
		for i, f := range frag.desc {
			bd.ui.Draw(gui.NewText(frag.pos[0], frag.pos[1]+i, f, &gui.TextConfig{
				FgColor: gui.White,
				BgColor: gui.Grey,
			}))
		}
	}
}

func (bd *BoardData) validateHitCoordinates(coord string) bool {
	x, y, err := mapCoords(coord)
	if err != nil || bd.app.OpponentBoardState[x][y] != gui.Empty {
		errorLog.SetText(fmt.Sprintf("%s", err))
		bd.statusAfterFire.SetText("Invalid coordinates, try again!")
		return false
	}
	bd.statusAfterFire.SetText("Valid coordinates.")
	return true
}

//func (a *App) handleShot() string {
//
//}

func (bd *BoardData) RenderBoards(status *models.StatusResponse) {
	hit := 0
	miss := 0

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(status.Timer)*time.Second)
	defer cancel()

	bd.playerTurn.SetText(fmt.Sprintf("Should I fire: %t", status.ShouldFire))
	bd.timer.SetText(fmt.Sprintf("Timer: %d", status.Timer))

	bd.drawBoard()

	//update timer
	go func() {
		for {
			status, _ = bd.app.client.GameStatus(GameStatusEndpoint)
			time.Sleep(time.Second / 2)
			bd.timer.SetText(fmt.Sprintf("Timer: %d", status.Timer))
			bd.playerTurn.SetText(fmt.Sprintf(If(status.ShouldFire,
				"It's your turn, fire at will", "It's your opponent turn, be patient")))
			_ = bd.markOpponentMoves(status)
		}
	}()

	//get input coords
	go func() {
		for {
			for status.ShouldFire {
				char := bd.opponentBoard.Listen(ctx)
				if bd.validateHitCoordinates(char) {

					shoot, _ := bd.app.client.Fire(FireEndpoint, char)
					var state gui.State

					if shoot.Result == "hit" || shoot.Result == "sunk" {
						bd.statusAfterFire.SetText(If(shoot.Result == "sunk", "Ship sunk", "Ship hit"))
						hit += 1
						state = gui.Hit
					} else {
						bd.statusAfterFire.SetText("")
						miss += 1
						state = gui.Miss
					}

					_ = bd.markPlayerMove(state, char)
					bd.playerMove.SetText(fmt.Sprintf("Last player move: %s", char))
				}
			}

		}
	}()

	// handle end game
	go func() {
		for {
			if status.GameStatus == "ended" {
				bd.gameResult.SetText(If(status.LastGameStatus == "win",
					"Game ended, You win", "Game ended, You lost"))
				bd.acc.SetText(fmt.Sprintf("Your accuracy: %.2f", float64(hit)/float64(miss)))
			}
		}
	}()

	bd.ui.Start(nil)
}
