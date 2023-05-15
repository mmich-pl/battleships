package main

import (
	. "battleships/internal/app"
	mainApp "battleships/internal/app"
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
			if err := StartNewGame(c); err != nil {
				log.Print(err)
				return
			}
		case "List players":
			if err := ListPlayer(c); err != nil {
				log.Println(err)
				continue
			}
		case "Print stats":
			if err := PrintTopTenPlayerStatistics(c); err != nil {
				log.Println(err)
				continue
			}
		case "Print stats of specific player":
			if err := PrintPlayerStatistics(c); err != nil {
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
				if err := StartNewGame(c); err != nil {
					log.Print(err)
					return
				}
			}
			gameInitiated = false
		}
	}

}
