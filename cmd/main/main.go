package main

import (
	. "battleships/internal/app"
	mainApp "battleships/internal/app/menu"
	"battleships/internal/battleship_client"

	"log"
)

var (
	baseURL       = "https://go-pjatk-server.fly.dev/api"
	gameInitiated = false
)

func main() {
	c := battleship_client.NewBattleshipClient(baseURL, 5, 5)
	menu := mainApp.InitializeMainMenu()

	for {
		playerChoide := menu.Display()
		switch playerChoide {
		case "Start game":
			if err := mainApp.StartNewGame(c); err != nil {
				log.Print(err)
				return
			}
		case "List players":
			if err := mainApp.ListPlayer(c); err != nil {
				log.Println(err)
				continue
			}
		case "Print stats":
			if err := mainApp.PrintTopTenPlayerStatistics(c); err != nil {
				log.Println(err)
				continue
			}
		case "Print stats of specific player":
			if err := mainApp.PrintPlayerStatistics(c); err != nil {
				log.Println(err)
				continue
			}
		case "Quit":
			return
		default:
			continue
		}
		if gameInitiated {
			input, _ := GetPlayerInput("play again? [yes]/[no]")
			if input == "yes" {
				if err := mainApp.StartNewGame(c); err != nil {
					log.Print(err)
					return
				}
			}
			gameInitiated = false
		}
	}

}
