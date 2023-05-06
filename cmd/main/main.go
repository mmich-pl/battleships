package main

import (
	mainApp "battleships/internal/app"
	"battleships/internal/battlehip_client"
	"log"
)

const (
	baseURL = "https://go-pjatk-server.fly.dev/api"
)

func main() {
	c := battlehip_client.NewBattleshipClient(baseURL, 5, 5)
	app := mainApp.New(c)
	err := app.Run()
	if err != nil {
		log.Print(err)
		return
	}
}
