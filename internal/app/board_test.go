package app

import (
	"battleships/internal/battlehip_client"
	"battleships/internal/models"
	_ "github.com/grupawp/warships-gui/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestMapCoords(t *testing.T) {
	testScenarios := []struct {
		name   string
		coords string
		error  bool
	}{
		{"valid coords - one int coords", "A1", false},
		{"valid coords - two int coords", "B10", false},
		{"invalid coords - two letter coords", "AB1", true},
		{"invalid coords - negative int", "A-1", true},
		{"invalid coords - letter out of bounds ", "Z1", true},
		{"invalid coords - int out of bounds", "A13", true},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			_, _, err := mapCoords(scenario.coords)
			if scenario.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRenderBoards(t *testing.T) {
	testScenario := struct {
		testName string
		client   func(t *testing.T) *battlehip_client.MockBattleshipClient
	}{
		testName: "Board Render",
		client: func(t *testing.T) *battlehip_client.MockBattleshipClient {
			client := battlehip_client.NewMockBattleshipClient(t)
			client.EXPECT().Description(OpponentDescription).Return(&models.DescriptionResponse{
				Desc:                "brytyjski admirał, urodzony w 1740 roku, zmarł w 1808 roku, dowódca sił morskich podczas wojen napoleońskich",
				Nick:                "Robert_Menzies",
				OpponentDescription: "Siejący trwogę, latający WP Bot. 999 walk wygranych przed czasem. Giń przeciwniku!",
				Opponent:            "WP_Bot",
			}, nil)

			client.EXPECT().GameStatus(GameStatusEndpoint).Return(&models.StatusResponse{
				GameStatus:     "game_in_porgress",
				LastGameStatus: "no_game",
				OpponentShots:  []string{"F4"},
				ShouldFire:     true,
				Timer:          36,
			}, nil)
			client.EXPECT().Board(BoardEndpoint).Return([]string{
				"A6", "A8", "A9", "C3", "D6", "D9", "D10", "E3", "F3", "F6",
				"G1", "G3", "G9", "G10", "H1", "H5", "I1", "J1", "J4", "J5"}, nil)
			return client
		},
	}

	t.Run(testScenario.testName, func(t *testing.T) {

		client := testScenario.client(t)
		app := New(client)
		status, _ := client.GameStatus(GameStatusEndpoint)
		board, _ := client.Board(BoardEndpoint)
		app.Description, _ = client.Description(OpponentDescription)

		if err := app.setUpBoardsState(board); err != nil {
			log.Fatal(err)
		}

		app.RenderBoards(status)
	})

}
