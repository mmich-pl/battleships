package app

import (
	"battleships/internal/battlehip_client"
	gui "github.com/grupawp/warships-gui"
	"github.com/stretchr/testify/assert"
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

func TestSetUpBoard(t *testing.T) {
	testScenario := struct {
		testName string
		client   func(t *testing.T) *battlehip_client.MockBattleshipClient
	}{
		testName: "Board Setup",
		client: func(t *testing.T) *battlehip_client.MockBattleshipClient {
			client := battlehip_client.NewMockBattleshipClient(t)
			client.EXPECT().Board(BoardEndpoint).Return([]string{
				"A7", "A8", "A9", "C3", "D6", "D9", "D10", "E3", "F3", "F6",
				"G1", "G3", "G9", "G10", "H1", "H5", "I1", "J1", "J4", "J5"}, nil)
			return client
		},
	}

	t.Run(testScenario.testName, func(t *testing.T) {
		client := testScenario.client(t)
		response, _ := client.Board(BoardEndpoint)
		expectedOpponent := [10][10]gui.State{}
		expectedPlayer := [10][10]gui.State{}
		for i := 0; i < len(expectedOpponent); i++ {
			expectedOpponent[i] = [10]gui.State{}
			expectedPlayer[i] = [10]gui.State{}
		}

		expectedPlayer[9][4] = gui.Ship
		expectedPlayer[0][7] = gui.Ship
		expectedPlayer[6][9] = gui.Ship
		expectedPlayer[3][5] = gui.Ship
		expectedPlayer[7][4] = gui.Ship
		expectedPlayer[9][3] = gui.Ship
		expectedPlayer[8][0] = gui.Ship
		expectedPlayer[3][8] = gui.Ship
		expectedPlayer[3][9] = gui.Ship
		expectedPlayer[2][2] = gui.Ship
		expectedPlayer[0][6] = gui.Ship
		expectedPlayer[4][2] = gui.Ship
		expectedPlayer[5][2] = gui.Ship
		expectedPlayer[9][0] = gui.Ship
		expectedPlayer[5][5] = gui.Ship
		expectedPlayer[6][0] = gui.Ship
		expectedPlayer[6][2] = gui.Ship
		expectedPlayer[6][8] = gui.Ship
		expectedPlayer[7][0] = gui.Ship
		expectedPlayer[0][8] = gui.Ship

		player, opponent, _ := setUpBoardsState(response)

		assert.Equalf(t, *player, expectedPlayer, "Expected to be the same")
		assert.Equalf(t, *opponent, expectedOpponent, "Expected to be the same")
	})
}
