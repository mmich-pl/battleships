package main

import (
	. "battleships/internal/app"
	. "battleships/internal/app/menu"
	"battleships/internal/battleship_client"
	. "battleships/internal/utils"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
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
	logFile := fmt.Sprintf("logs/log-%v.txt", time.Now().Format(time.RFC3339))
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + logFile)
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	err = godotenv.Load()
	if err != nil {
		return
	}
	c := battleship_client.NewBattleshipClient(baseURL, 5, 5)
	app := New(c)
	menu := InitializeMainMenu()

	for {
		playerChoide := menu.Display()
		switch playerChoide {
		case "Start game":
			if err = startNewGame(app); err != nil {
				log.Print(err)
				return
			}
			gameInitiated = true
		case "List players":
			if err = ListPlayer(c); err != nil {
				log.Println(err)
				continue
			}
		case "Print stats":
			if err = PrintTopTenPlayerStatistics(c); err != nil {
				log.Println(err)
				continue
			}
		case "Print stats of specific player":
			if err = PrintPlayerStatistics(c); err != nil {
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
