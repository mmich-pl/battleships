package menu

import (
	"battleships/internal/battleship_client"
	models "battleships/internal/models"
	"testing"
)

func TestMenu_Display(t *testing.T) {
	testScenario := struct {
		testName string
		client   func(t *testing.T) *battleship_client.MockBattleshipClient
	}{
		testName: "Board Render",
		client: func(t *testing.T) *battleship_client.MockBattleshipClient {
			client := battleship_client.NewMockBattleshipClient(t)
			client.EXPECT().GetPlayersList().Return(
				&[]models.WaitingPlayerData{
					{
						GameStatus: "waiting",
						Nick:       "Kapitan1",
					},
					{
						GameStatus: "waiting",
						Nick:       "Kapitan2",
					},
					{
						GameStatus: "waiting",
						Nick:       "Kapitan4",
					},
					{
						GameStatus: "waiting",
						Nick:       "Kapitan3",
					},
				}, nil)
			return client
		},
	}

	t.Run(testScenario.testName, func(t *testing.T) {

		client := testScenario.client(t)
		menu := NewMenu("Welcome in Battleship, what do you want to do?")
		playerList, _ := client.GetPlayersList()

		for _, player := range *playerList {
			menu.AddItem(player.Nick)
		}

		menu.Display()
	})

}
