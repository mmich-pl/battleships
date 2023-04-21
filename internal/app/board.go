package app

import (
	"fmt"
	gui "github.com/grupawp/warships-gui"
	"golang.org/x/sync/errgroup"
	"strconv"
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
		return -1, -1, fmt.Errorf("coordinate out of bound: expected row in bounds [1-10]")
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
		g.Go(func() error {
			if x, y, err := mapCoords(coords); err != nil {
				return err
			} else {
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
