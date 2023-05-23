package app

import (
	. "battleships/internal/app/board_validation"
	"battleships/internal/utils"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"time"
)

func getPlayerFleet(stop chan struct{}, states [10][10]gui.State, board gui.Board, txt gui.Text, getPlayerFleet chan<- []string) {
	playerFleet := make([]string, 0)

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
			getPlayerFleet <- playerFleet
			return
		default:
			coords := board.Listen(context.TODO())
			x, y, _ := utils.MapCoords(coords)
			switch s := &states[x][y]; *s {
			case gui.Ship:
				*s = gui.Empty
				playerFleet = remove(playerFleet, coords)
				txt.SetText(fmt.Sprintf("Ship remove from %s", coords))
			default:
				*s = gui.Ship
				playerFleet = append(playerFleet, coords)
				txt.SetText(fmt.Sprintf("Ship set as %s", coords))
			}
			board.SetStates(states)

			if len(playerFleet) >= 20 {
				board.SetStates(states)
				txt.SetText("All fleet set!")
				time.Sleep(time.Second * 2)
				getPlayerFleet <- playerFleet
				return
			}
		}
	}
}

func printFleetInstruction(ui gui.GUI) {
	title := gui.NewText(60, 10, "Your fleet should look like this:", nil)

	battleshipInfo := gui.NewText(60, 11, "1x Battleship  (4 tiles)", nil)
	cruiserInfo := gui.NewText(60, 12, "2x Cruiser     (3 tiles)", nil)
	destroyerInfo := gui.NewText(60, 13, "3x Destroyer   (2 tiles)", nil)
	submarineInfo := gui.NewText(60, 14, "4x Submarine   (1 tile)", nil)

	ui.Draw(title)
	ui.Draw(battleshipInfo)
	ui.Draw(cruiserInfo)
	ui.Draw(destroyerInfo)
	ui.Draw(submarineInfo)
}

func initBaseState() [10][10]gui.State {
	states := [10][10]gui.State{}
	for i := range states {
		states[i] = [10]gui.State{}
	}
	return states
}

func RenderInputBoard() []string {
	var playerFleet []string
	board := gui.NewBoard(1, 1, nil)
	ui := gui.NewGUI(false)
	states := initBaseState()
	board.SetStates(states)
	ui.Draw(board)

	txt := gui.NewText(60, 5, "Press on any coordinate to set a ship.", nil)
	ui.Draw(txt)
	printFleetInstruction(*ui)

	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan struct{})
	fleetChannel := make(chan []string)

	go func() {
		validBoard := false
		var cause string
		for !validBoard {
			go getPlayerFleet(stop, states, *board, *txt, fleetChannel)
			playerFleet = <-fleetChannel
			validBoard, cause = ValidateShipPlacement(playerFleet)
			txt.SetText(cause)
			states = initBaseState()
			board.SetStates(states)
		}
		close(stop)
		cancel()
	}()
	ui.Start(ctx, nil)

	return playerFleet
}
