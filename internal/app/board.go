package app

import (
	"battleships/internal/models"
	. "battleships/internal/utils"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

var (
	playerDescPos   = [2]int{2, 29}
	opponentDescPos = [2]int{50, 29}
)

type BoardData struct {
	app                    *App
	ui                     *gui.GUI
	playerBoardIndicator   *gui.Text
	opponentBoardIndicator *gui.Text
	playerBoard            *gui.Board
	opponentBoard          *gui.Board
	playerNick             *gui.Text
	opponentNick           *gui.Text
	playerTurn             *gui.Text
	timer                  *gui.Text
	gameResult             *gui.Text
	statusAfterFire        *gui.Text
	playerMove             *gui.Text
	accuracy               *gui.Text
	playerStatsHeader      *gui.Text
	playerStats            *gui.Text
	opponentStatsHeader    *gui.Text
	opponentStats          *gui.Text
	legend                 *gui.Text
	instructions           *gui.Text
}

func InitBoardData(a *App) *BoardData {
	return &BoardData{
		app:                    a,
		ui:                     gui.NewGUI(false),
		playerBoardIndicator:   gui.NewText(2, 5, "Player board", nil),
		opponentBoardIndicator: gui.NewText(65, 5, "Opponent board", nil),
		playerBoard:            gui.NewBoard(2, 6, nil),
		opponentBoard:          gui.NewBoard(50, 6, nil),
		playerNick:             gui.NewText(2, 28, a.Description.Nick, nil),
		opponentNick:           gui.NewText(50, 28, a.Description.Opponent, nil),
		playerTurn:             gui.NewText(65, 3, "", nil),
		timer:                  gui.NewText(50, 3, "", nil),
		gameResult:             gui.NewText(2, 32, fmt.Sprintf("Game is running!"), nil),
		statusAfterFire:        gui.NewText(2, 31, "", nil),
		playerMove:             gui.NewText(2, 3, "Press on any coordinate to take a shot.", nil),
		accuracy:               gui.NewText(2, 33, "", nil),
		playerStatsHeader:      gui.NewText(97, 5, "", nil),
		playerStats:            gui.NewText(97, 6, "", nil),
		opponentStatsHeader:    gui.NewText(97, 8, "", nil),
		opponentStats:          gui.NewText(97, 9, "", nil),
		legend:                 gui.NewText(97, 12, "Symbols: S - indicate ship, M - indicate miss, H - indicate hit", nil),
		instructions:           gui.NewText(97, 13, "To perform hit you have to click box on opponent board. It will work only if it is your turn.", nil),
	}
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

func (bd *BoardData) markPlayerMove(state gui.State, coord string) error {
	x, y, err := MapCoords(coord)
	if err != nil {
		return fmt.Errorf("failed to parse coord: %w", err)
	}

	bd.app.OpponentBoardState[x][y] = state
	bd.opponentBoard.SetStates(bd.app.OpponentBoardState)
	return nil
}

func (bd *BoardData) markOpponentMoves(status *models.StatusResponse) error {

	for _, cords := range status.OpponentShots {
		x, y, err := MapCoords(cords)
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
	bd.ui.Draw(bd.playerBoardIndicator)
	bd.ui.Draw(bd.opponentBoard)
	bd.ui.Draw(bd.opponentBoardIndicator)

	bd.ui.Draw(bd.playerTurn)
	bd.ui.Draw(bd.timer)
	bd.ui.Draw(bd.gameResult)
	bd.ui.Draw(bd.statusAfterFire)
	bd.ui.Draw(bd.playerMove)
	bd.ui.Draw(bd.accuracy)

	opponent, _ := bd.app.Client.GetPlayerStatistic(bd.app.Description.Opponent)
	player, _ := bd.app.Client.GetPlayerStatistic(bd.app.Description.Nick)

	bd.playerStatsHeader.SetText(fmt.Sprintf("%7s |\t%20s |\t%6s |\t%6s|\t%6s", "Rank", "Nick", "Points", "Games", "Wins"))
	bd.opponentStatsHeader.SetText(fmt.Sprintf("%7s |\t%20s |\t%6s |\t%6s|\t%6s", "Rank", "Nick", "Points", "Games", "Wins"))
	bd.playerStats.SetText(fmt.Sprintf("%7d |\t%20s |\t%6d |\t%6d|\t%6d",
		player.Stats.Rank, player.Stats.Nick, player.Stats.Points, player.Stats.Games, player.Stats.Wins))
	bd.opponentStats.SetText(fmt.Sprintf("%7d |\t%20s |\t%6d |\t%6d|\t%6d",
		opponent.Stats.Rank, opponent.Stats.Nick, opponent.Stats.Points, opponent.Stats.Games, opponent.Stats.Wins))
	bd.ui.Draw(bd.playerStatsHeader)
	bd.ui.Draw(bd.playerStats)

	bd.ui.Draw(bd.opponentStatsHeader)
	bd.ui.Draw(bd.opponentStats)
	bd.ui.Draw(bd.legend)
	bd.ui.Draw(bd.instructions)
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

func (bd *BoardData) handleShot() string {
	for {
		coords := bd.opponentBoard.Listen(context.TODO())
		x, y, _ := MapCoords(coords)
		if bd.app.OpponentBoardState[x][y] == gui.Hit || bd.app.OpponentBoardState[x][y] == gui.Miss {
			bd.statusAfterFire.SetText("Invalid coordinates, try again!")
		} else {
			bd.statusAfterFire.SetText("Valid coordinates.")
			bd.playerMove.SetText(fmt.Sprintf("Last player move: %s", coords))
			return coords
		}
	}
}

func (bd *BoardData) RenderGameBoards(status *models.StatusResponse) error {
	hit := 0
	miss := 0

	bd.playerTurn.SetText(fmt.Sprintf("Should I fire: %t", status.ShouldFire))
	bd.timer.SetText(fmt.Sprintf("Timer: %d", status.Timer))

	bd.drawBoard()
	go func() {
		for status.GameStatus == "game_in_progress" {
			status, _ = bd.app.Client.GameStatus()
			bd.timer.SetText(fmt.Sprintf("Timer: %d", status.Timer))
			time.Sleep(time.Second)
		}
	}()

	//game logic
	go func() {
		for status.GameStatus == "game_in_progress" {

			time.Sleep(time.Second)
			_ = bd.markOpponentMoves(status)
			bd.playerTurn.SetText(fmt.Sprintf(If(status.ShouldFire,
				"It's your turn, fire at will", "It's your opponent turn, be patient")))

			shouldContinue := true

			for shouldContinue && status.ShouldFire && status.GameStatus != "ended" {
				coords := bd.handleShot()
				if len(coords) != 0 {
					shoot, _ := bd.app.Client.Fire(coords)

					var state gui.State
					if shoot.Result == "hit" || shoot.Result == "sunk" {
						bd.statusAfterFire.SetText(If(shoot.Result == "sunk", "Ship sunk", "Ship hit"))
						state = gui.Hit
						hit++
					} else {
						bd.statusAfterFire.SetText("")
						state = gui.Miss
						shouldContinue = false
						miss++
					}

					_ = bd.markPlayerMove(state, coords)
				}
			}

		}
	}()
	boardCtx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			if status.GameStatus == "ended" {
				bd.gameResult.SetText(If(status.LastGameStatus == "win",
					"Game ended, You win", "Game ended, You lost"))
				bd.accuracy.SetText(fmt.Sprintf("%.2f", If(hit != 0 && miss != 0, float64(hit)/float64(miss+hit), 0)))
				time.Sleep(5 * time.Second)
				cancel()
			}
		}

	}()

	bd.ui.Start(boardCtx, nil)
	err := bd.app.Client.AbandonGame()
	if err != nil {
		return err
	}
	return nil
}
