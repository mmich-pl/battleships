package app

import (
	"battleships/internal/models"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui"
	"golang.org/x/sync/errgroup"
	"log"
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

func RenderBoards(status *models.FullStatusResponse, playerState, opponentState *[10][10]gui.State) error {
	ctx := context.TODO()

	drawer := gui.NewDrawer(&gui.Config{})

	playerBoard, err := drawer.NewBoard(2, 4, &gui.BoardConfig{})
	if err != nil {
		return fmt.Errorf("failed to render player board: %w", err)
	}

	opponentNick, err := drawer.NewText(2, 1, nil)
	if err != nil {
		return fmt.Errorf("failed to render player nick: %w", err)
	}

	opponentDescription, err := drawer.NewText(2, 3, nil)
	if err != nil {
		return fmt.Errorf("failed to render player description: %w", err)
	}

	opponentBoard, err := drawer.NewBoard(50, 4, &gui.BoardConfig{})
	if err != nil {
		return fmt.Errorf("failed to render opponent board: %w", err)
	}

	defer drawer.RemoveBoard(ctx, playerBoard)
	defer drawer.RemoveBoard(ctx, opponentBoard)
	defer drawer.RemoveText(ctx, opponentNick)
	defer drawer.RemoveText(ctx, opponentDescription)

	drawer.DrawBoard(ctx, playerBoard, *playerState)
	drawer.DrawBoard(ctx, opponentBoard, *opponentState)

	opponentNick.SetText(status.Opponent)
	opponentDescription.SetText(status.OpponentDescription)

	drawer.DrawText(ctx, opponentNick)
	drawer.DrawText(ctx, opponentDescription)

	for {
		if !drawer.IsGameRunning() {
			return nil
		}
	}
}
