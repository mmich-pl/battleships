package app

import (
	"battleships/internal/app/board_utils"
	"battleships/internal/models"
	. "battleships/internal/utils"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/mitchellh/go-wordwrap"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

var (
	ShipNames         = []string{"Submarine", "Destroyer", "Cruiser", "Battleship"}
	playerDescPos     = [2]int{115, 7}
	opponentDescPos   = [2]int{170, 7}
	playerBoardConfig = &gui.BoardConfig{
		TextColor:  gui.Color{},
		RulerColor: gui.Color{Red: 255, Green: 255, Blue: 255},
		EmptyColor: gui.Color{Red: 155, Green: 234, Blue: 236},
		HitColor:   gui.Color{Red: 251, Green: 5, Blue: 74},
		MissColor:  gui.Color{Red: 81, Green: 132, Blue: 237},
		ShipColor:  gui.Color{Red: 133, Green: 133, Blue: 133},
		EmptyChar:  '~',
		HitChar:    'H',
		MissChar:   'M',
		ShipChar:   'S',
	}
	opponentBoardConfig = &gui.BoardConfig{
		TextColor:  gui.Color{},
		RulerColor: gui.Color{Red: 255, Green: 255, Blue: 255},
		EmptyColor: gui.Color{Red: 202, Green: 202, Blue: 216},
		HitColor:   gui.Color{Red: 251, Green: 5, Blue: 74},
		MissColor:  gui.Color{Red: 133, Green: 133, Blue: 133},
		EmptyChar:  '~',
		HitChar:    'H',
		MissChar:   'M',
		ShipChar:   'S',
	}
)

type BoardData struct {
	app                    *App
	ui                     *gui.GUI
	playerBoardIndicator   *gui.Text
	opponentBoardIndicator *gui.Text
	playerBoard            *gui.Board
	opponentBoard          *gui.Board
	opponentFleet          map[int]int
	opponentFleetTable     []*gui.Text
	playerNick             *gui.Text
	opponentNick           *gui.Text
	playerTurn             *gui.Text
	timer                  *gui.Text
	gameResult             *gui.Text
	statusAfterFire        *gui.Text
	playerMove             *gui.Text
	accuracy               *gui.Text
	legend                 *gui.Text
	instructions           *gui.Text
}

func InitBoardData(a *App) *BoardData {
	opponentFleet := make(map[int]int)
	playerFleet = make(map[int]int)
	opponentTable := make([]*gui.Text, 5)

	opponentTable[0] = gui.NewText(70, 29, "", nil)

	for k, v := range board_utils.ShipQuantities {
		opponentFleet[k] = v
		playerFleet[k] = v
		opponentTable[k] = gui.NewText(70, 29+k, "", nil)
	}

	return &BoardData{
		app:                    a,
		ui:                     gui.NewGUI(false),
		playerBoardIndicator:   gui.NewText(2, 5, "Player board_utils", nil),
		opponentBoardIndicator: gui.NewText(72, 5, "Opponent board_utils", nil),
		playerBoard:            gui.NewBoard(2, 6, playerBoardConfig),
		opponentBoard:          gui.NewBoard(70, 6, opponentBoardConfig),
		opponentFleet:          opponentFleet,
		opponentFleetTable:     opponentTable,
		playerNick:             gui.NewText(115, 6, a.Description.Nick, nil),
		opponentNick:           gui.NewText(170, 6, a.Description.Opponent, nil),
		playerTurn:             gui.NewText(65, 3, "", nil),
		timer:                  gui.NewText(50, 3, "", nil),
		gameResult:             gui.NewText(115, 3, fmt.Sprintf("Game is running!"), nil),
		statusAfterFire:        gui.NewText(135, 3, "", nil),
		playerMove:             gui.NewText(2, 3, "Press on any coordinate to take a shot.", nil),
		accuracy:               gui.NewText(150, 3, "", nil),
		legend:                 gui.NewText(115, 20, "Symbols: S - indicate ship, M - indicate miss, H - indicate hit", nil),
		instructions:           gui.NewText(115, 21, "To perform hit you have to click box on opponent board_utils. It will work only if it is your turn.", nil),
	}
}

// Parses API response to [10][10] matrix format used by client
func (a *App) setUpBoardsState(board []string) error {
	for i := 0; i < board_utils.BoardSize; i++ {
		a.PlayerBoardState[i] = [10]gui.State{}
		a.OpponentBoardState[i] = [10]gui.State{}
	}

	g := new(errgroup.Group)
	for _, coords := range board {
		c := coords
		g.Go(func() error {
			if x, y, err := MapCoords(c); err != nil {
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

func (bd *BoardData) markPlayerMove(state gui.State, x, y int, result string) error {
	bd.app.OpponentBoardState[x][y] = state
	if result == "sunk" {
		bd.markBorder(x, y)
	}
	bd.opponentBoard.SetStates(bd.app.OpponentBoardState)
	return nil
}

func (bd *BoardData) markOpponentMoves(status *models.StatusResponse) error {
	var x, y int
	var err error
	for _, cords := range status.OpponentShots {
		x, y, err = MapCoords(cords)
		if err != nil {
			log.Error(err)
			return fmt.Errorf("failed to parse coords: %w", err)
		}

		switch state := &bd.app.PlayerBoardState[x][y]; *state {
		case gui.Hit, gui.Ship:
			*state = gui.Hit
		default:
			*state = gui.Miss
		}
	}
	return nil
}

func (bd *BoardData) printFleetInfo(table []*gui.Text) {
	table[0].SetText(fmt.Sprintf("%12s |\t%6s |\t%16s |\t%15s", "Ship", "Size", "Initial amount", "Survived ships"))
	bd.ui.Draw(table[0])
	for k, v := range bd.opponentFleet {
		table[k].SetText(fmt.Sprintf("%12s |\t%6d |\t%16d |\t%15d", ShipNames[k-1], k, board_utils.ShipQuantities[k], v))
		bd.ui.Draw(table[k])
	}
}

func (bd *BoardData) printStats(x, y int) {
	opponent, err := bd.app.Client.GetPlayerStatistic(bd.app.Description.Opponent)
	if err != nil {
		log.Error(err)
	}
	player, err := bd.app.Client.GetPlayerStatistic(bd.app.Description.Nick)
	if err != nil {
		log.Error(err)
	}
	temp := []*gui.Text{
		gui.NewText(x, y, fmt.Sprintf("|\t%10s|\t%20s|\t%20s|", "Statistic", player.Stats.Nick, opponent.Stats.Nick), nil),
		gui.NewText(x, y+1, fmt.Sprintf("|\t%10s|\t%20d|\t%20d|", "Rank", player.Stats.Rank, opponent.Stats.Rank), nil),
		gui.NewText(x, y+2, fmt.Sprintf("|\t%10s|\t%20d|\t%20d|", "Points", player.Stats.Points, opponent.Stats.Points), nil),
		gui.NewText(x, y+3, fmt.Sprintf("|\t%10s|\t%20d|\t%20d|", "Games", player.Stats.Games, opponent.Stats.Games), nil),
		gui.NewText(x, y+4, fmt.Sprintf("|\t%10s|\t%20d|\t%20d|", "Wins", player.Stats.Wins, opponent.Stats.Wins), nil),
	}

	for _, i := range temp {
		bd.ui.Draw(i)
	}
}

func (bd *BoardData) drawBoard() {
	bd.playerBoard.SetStates(bd.app.PlayerBoardState)
	bd.opponentBoard.SetStates(bd.app.OpponentBoardState)

	bd.ui.Draw(bd.playerNick)
	bd.ui.Draw(bd.opponentNick)

	bd.renderDescription()

	bd.ui.Draw(bd.playerBoard)
	bd.ui.Draw(bd.playerBoardIndicator)
	bd.ui.Draw(bd.opponentBoard)
	bd.ui.Draw(bd.opponentBoardIndicator)

	bd.ui.Draw(bd.playerTurn)
	bd.ui.Draw(bd.timer)
	bd.ui.Draw(bd.gameResult)
	bd.ui.Draw(bd.statusAfterFire)
	bd.ui.Draw(bd.playerMove)
	bd.ui.Draw(bd.accuracy)

	bd.ui.Draw(bd.legend)
	bd.ui.Draw(bd.instructions)
	bd.printStats(120, 12)
	bd.printFleetInfo(bd.opponentFleetTable)
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
		if len(frag.desc) == 0 {
			continue
		}

		for i, f := range frag.desc {
			bd.ui.Draw(gui.NewText(frag.pos[0], frag.pos[1]+i, f, &gui.TextConfig{
				FgColor: gui.White,
				BgColor: gui.Grey,
			}))
		}
	}

}

func (bd *BoardData) validateShootCoordinates() string {
	for {
		coords := bd.opponentBoard.Listen(context.TODO())
		x, y, err := MapCoords(coords)
		if err != nil {
			log.Error(err)
		}
		if bd.app.OpponentBoardState[x][y] == gui.Hit || bd.app.OpponentBoardState[x][y] == gui.Miss {
			bd.statusAfterFire.SetText("Invalid coordinates, try again!")
		} else {
			bd.statusAfterFire.SetText("Valid coordinates.")
			bd.playerMove.SetText(fmt.Sprintf("Last player move: %s", coords))
			return coords
		}
	}
}

func (bd *BoardData) handleShoot(coords string, err error, hit int, shouldContinue bool, miss int) (int, int, bool) {
	if len(coords) != 0 {
		var x, y int
		var shoot *models.ShootResult
		x, y, err = MapCoords(coords)
		if err != nil {
			log.Error(err)
		}
		shoot, err = bd.app.Client.Fire(coords)
		if err != nil {
			log.Error(err)
		}
		var state gui.State

		switch shoot.Result {
		case "hit", "sunk":
			bd.statusAfterFire.SetText("Ship hit")
			state = gui.Hit
			hit++
		default:
			bd.statusAfterFire.SetText("You miss")
			state = gui.Miss
			shouldContinue = false
			miss++
		}

		err = bd.markPlayerMove(state, x, y, shoot.Result)
		if err != nil {
			log.Error(err)
		}
	}
	return hit, miss, shouldContinue
}

func (bd *BoardData) RenderGameBoards(status *models.StatusResponse) error {
	hit := 0
	miss := 0

	bd.playerTurn.SetText(fmt.Sprintf("Should I fire: %t", status.ShouldFire))
	bd.timer.SetText(fmt.Sprintf("Timer: %d", status.Timer))

	bd.drawBoard()
	go func() {
		var err error
		for status.GameStatus == "game_in_progress" {
			status, err = bd.app.Client.GameStatus()
			if err != nil {
				log.Error(err)
			}
			bd.timer.SetText(fmt.Sprintf("Timer: %d", status.Timer))
			time.Sleep(time.Second)
		}
	}()

	//game logic
	go func() {
		for status.GameStatus == "game_in_progress" {

			time.Sleep(time.Second)
			err := bd.markOpponentMoves(status)
			if err != nil {
				log.Error(err)
			}
			bd.playerTurn.SetText(fmt.Sprintf(If(status.ShouldFire,
				"It's your turn, fire at will", "It's your opponent turn, be patient")))

			shouldContinue := true

			for shouldContinue && status.ShouldFire && status.GameStatus != "ended" {
				coords := bd.validateShootCoordinates()
				hit, miss, shouldContinue = bd.handleShoot(coords, err, hit, shouldContinue, miss)
			}

		}
	}()
	boardCtx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			if status.GameStatus == "ended" {
				bd.gameResult.SetText(If(status.LastGameStatus == "win",
					"Game ended, You win", "Game ended, You lost"))
				bd.accuracy.SetText(fmt.Sprintf("Your accuracy: %.2f%% (%d/%d)", If(hit != 0 && miss != 0,
					(float64(hit)/float64(miss+hit))*100, 0), hit, miss+hit))
				time.Sleep(5 * time.Second)
				cancel()
			}
		}

	}()

	bd.ui.Start(boardCtx, nil)
	if status.GameStatus != "ended" {
		err := bd.app.Client.AbandonGame()
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}
