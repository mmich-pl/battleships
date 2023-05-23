package main

import (
	. "battleships/internal/app"
	mainApp "battleships/internal/app/menu"
	"battleships/internal/battleship_client"
	. "battleships/internal/utils"
	"github.com/joho/godotenv"
	"log"
)

var (
	baseURL       = "https://go-pjatk-server.fly.dev/api"
	gameInitiated = false
)

func startNewGame(app *App) error {
	err := app.Run()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	c := battleship_client.NewBattleshipClient(baseURL, 5, 5)
	app := New(c)
	menu := mainApp.InitializeMainMenu()

	for {
		playerChoide := menu.Display()
		switch playerChoide {
		case "Start game":
			if err := startNewGame(app); err != nil {
				log.Print(err)
				return
			}
			gameInitiated = true
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
			input, _ := GetPlayerInput("play again? [y/n]: ")
			if input == "y" {
				if err := startNewGame(app); err != nil {
					log.Print(err)
					return
				}
			}
			gameInitiated = false
		}
	}

}
