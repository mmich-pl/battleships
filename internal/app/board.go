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
)

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

func setUpBoardsState(board []string) (*[10][10]gui.State, *[10][10]gui.State, error) {
	var playerBoardState, opponentBoardState [10][10]gui.State
	for i := 0; i < len(playerBoardState); i++ {
		playerBoardState[i] = [10]gui.State{}
		opponentBoardState[i] = [10]gui.State{}
	}

	g := new(errgroup.Group)
	for _, coords := range board {
		c := coords
		g.Go(func() error {
			if x, y, err := mapCoords(c); err != nil {
				return err
			} else {
				log.Printf("set ship om position: [%d, %d]", x, y)
				playerBoardState[x][y] = gui.Ship
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, fmt.Errorf("failed to fill boards: %w", err)
	}

	return &playerBoardState, &opponentBoardState, nil
}

func RenderDescription(g gui.GUI, playerDescription, opponentDescription string) {
	ypos := 25

	fragments := [2]struct {
		desc []string
		xpos int
	}{
		{strings.Split(wordwrap.WrapString(playerDescription, 40), "\n"), 2},
		{strings.Split(wordwrap.WrapString(opponentDescription, 40), "\n"), 50},
	}
	for _, frag := range fragments {
		for i, f := range frag.desc {
			g.Draw(gui.NewText(frag.xpos, ypos+i, f, &gui.TextConfig{
				FgColor: gui.White,
				BgColor: gui.Grey,
			}))
		}
	}

}

func RenderBoards(status *models.FullStatusResponse, playerState, opponentState [10][10]gui.State) {
	ui := gui.NewGUI(true)
	playerBoard := gui.NewBoard(2, 3, nil)
	playerBoard.SetStates(playerState)
	opponentBoard := gui.NewBoard(50, 3, nil)
	opponentBoard.SetStates(opponentState)

	ui.Draw(gui.NewText(2, 1, fmt.Sprintf("%s vs %s", status.Nick, status.Opponent), nil))
	RenderDescription(*ui, status.Desc, status.OpponentDescription)

	ui.Draw(playerBoard)
	ui.Draw(opponentBoard)

	playerMove := gui.NewText(2, 29, "Press on any coordinate to log it.", nil)
	ui.Draw(playerMove)

	go func() {
		for {
			char := opponentBoard.Listen(context.TODO())
			playerMove.SetText(fmt.Sprintf("Ready! Aim at %v! FIRE!", char))
			ui.Log("Coordinate: %s", char) // logs are displayed after the game exits
		}
	}()

	ui.Start(nil)

}
