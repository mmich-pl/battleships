package main

import (
	"battleships/internal/battlehip_client"
	"log"
)

const (
	baseURL      = "https://go-pjatk-server.fly.dev/api"
	InitEndpoint = "/game"
)

func main() {
	c := battlehip_client.NewBattleshipClient(baseURL, 5, 5)
	err := c.InitGame(InitEndpoint, "", "", "", false)
	if err != nil {
		log.Print(err)
		return
	}

}
