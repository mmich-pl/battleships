package main

import (
	mainApp "battleships/internal/app"
	"battleships/internal/battleship_client"
	"fmt"
	"log"
	"strings"
)

const (
	baseURL = "https://go-pjatk-server.fly.dev/api"
)

func main() {
	c := battleship_client.NewBattleshipClient(baseURL, 5, 5)
	//app := mainApp.New(c)
	menu := mainApp.InitializeMainMenu()

	for {
		playerChoide := menu.Display()
		switch playerChoide {
		case "List players":
			playerList, err := c.GetPlayersList(mainApp.WaitingPlayersEndpoint)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, data := range *playerList {
				fmt.Println(data.GameStatus, " | ", data.Nick)
			}

		case "Print stats":
			stats, err := c.GetStatistic(mainApp.StatsEndpoint)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, s := range stats.Stats {
				fmt.Println("Nick:", s.Nick, "Games:", s.Games,
					"Wins:", s.Wins, "Rank:", s.Rank, "Points", s.Points)
			}
		case "Print stats of specific player":
			input, _ := mainApp.GetPlayerInput("get the stats for player:")

			s, err := c.GetPlayerStatistic(mainApp.StatsEndpoint, strings.TrimSpace(input))
			if err != nil {
				log.Println(err)
				continue
			}

			fmt.Println("Nick:", s.Stats.Nick, "Games:", s.Stats.Games,
				"Wins:", s.Stats.Wins, "Rank:", s.Stats.Rank, "Points", s.Stats.Points)

		case "Quit":
			return
		default:
			continue
		}

	}

	//err := app.Run()
	//if err != nil {
	//	log.Print(err)
	//	return
	//}
}
