package app

import (
	"battleships/internal/app/board_utils"
	"battleships/internal/utils"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"time"
)

var (
	ui          = gui.NewGUI(false)
	playerFleet = make(map[int]int)
)

func getPlayerFleet(stop chan struct{}, states [10][10]gui.State, board gui.Board, txt gui.Text, getPlayerFleet chan<- []string) {
	playerFleetCoords := make([]string, 0)

	remove := func(s []string, element string) []string {
		var index int
		for i, e := range s {
			if element == e {
				index = i
			}
		}
		return append(s[:index], s[index+1:]...)
	}

	for {
		select {
		case <-stop:
			getPlayerFleet <- playerFleetCoords
			return
		default:
			coords := board.Listen(context.TODO())
			x, y, _ := utils.MapCoords(coords)
			switch s := &states[x][y]; *s {
			case gui.Ship:
				*s = gui.Empty
				playerFleetCoords = remove(playerFleetCoords, coords)
				txt.SetText(fmt.Sprintf("Ship remove from %s", coords))
			default:
				*s = gui.Ship
				playerFleetCoords = append(playerFleetCoords, coords)
				txt.SetText(fmt.Sprintf("Ship set at %s (%d)", coords, len(coords)))
			}
			updatePlayerFleet(playerFleetCoords)
			board.SetStates(states)

			if len(playerFleetCoords) == 20 {
				board.SetStates(states)
				txt.SetText("All fleet set!")
				time.Sleep(time.Second * 2)
				getPlayerFleet <- playerFleetCoords
				return
			}
		}
	}
}

func updatePlayerFleet(playerFleetCoords []string) {
	b, _ := board_utils.MapCoordsToBoard(playerFleetCoords)
	blobs := board_utils.ConnectedComponentLabeling(b)
	playerFleet = board_utils.CountOccurrences(blobs)
	printFleetInstruction()
}

func printFleetInstruction() {
	header := gui.NewText(60, 10, fmt.Sprintf("%12s |\t%6s |\t%16s ", "Ship", "Size", "Can be placed"), nil)
	ui.Draw(header)

	for k, v := range board_utils.ShipQuantities {
		t := gui.NewText(60, 10+k, fmt.Sprintf("%12s |\t%6d |\t%16d ", ShipNames[k-1], k, v-playerFleet[k]), nil)
		ui.Draw(t)
	}
}

func initBaseState() [10][10]gui.State {
	states := [10][10]gui.State{}
	for i := range states {
		states[i] = [10]gui.State{}
	}
	return states
}

func RenderInputBoard() []string {
	for k, _ := range board_utils.ShipQuantities {
		playerFleet[k] = 0
	}
	var playerFleetCoords []string
	board := gui.NewBoard(1, 1, nil)

	states := initBaseState()
	board.SetStates(states)
	ui.Draw(board)

	txt := gui.NewText(60, 5, "Press on any coordinate to set a ship.", nil)
	ui.Draw(txt)
	printFleetInstruction()

	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan struct{})
	fleetChannel := make(chan []string)

	go func() {
		validBoard := false
		var cause string
		for !validBoard {
			go getPlayerFleet(stop, states, *board, *txt, fleetChannel)
			playerFleetCoords = <-fleetChannel
			validBoard, cause = board_utils.ValidateShipPlacement(playerFleetCoords)
			txt.SetText(cause)
			states = initBaseState()
			board.SetStates(states)
		}
		close(stop)
		cancel()
	}()
	ui.Start(ctx, nil)

	return playerFleetCoords
}
