package app

import (
	"battleships/internal/battleship_client"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

func InitializeMainMenu() *Menu {
	menu := NewMenu("Welcome in Battleship, what do you want to do?")
	menu.AddItem("List players")
	menu.AddItem("Print stats")
	menu.AddItem("Print stats of specific player")
	menu.AddItem("Start game")
	menu.AddItem("Quit")
	return menu
}

func PrintPlayerStatistics(c *battleship_client.BattleshipHTTPClient) error {
	input, _ := GetPlayerInput("get the stats for player:")

	s, err := c.GetPlayerStatistic(StatsEndpoint, strings.TrimSpace(input))
	if err != nil {
		log.Println(err)
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', tabwriter.AlignRight)

	_, _ = fmt.Fprintf(w, "%s |\t%s |\t%s |\t%s|\t%s\n", "Rank", "Nick", "Points", "Games", "Wins")
	_, _ = fmt.Fprintf(w, "%d |\t%s |\t%d |\t%d|\t%d\n",
		s.Stats.Rank, s.Stats.Nick, s.Stats.Points, s.Stats.Games, s.Stats.Wins)
	_ = w.Flush()

	return nil
}

func PrintTopTenPlayerStatistics(c *battleship_client.BattleshipHTTPClient) error {
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', tabwriter.AlignRight)
	stats, err := c.GetStatistic(StatsEndpoint)
	if err != nil {

		return err
	}

	_, _ = fmt.Fprintf(w, "%s |\t%s |\t%s |\t%s|\t%s\n", "Rank", "Nick", "Points", "Games", "Wins")
	for _, s := range stats.Stats {
		_, _ = fmt.Fprintf(w, "%d |\t%s |\t%d |\t%d|\t%d\n", s.Rank, s.Nick, s.Points, s.Games, s.Wins)
	}
	_ = w.Flush()
	return nil
}

func ListPlayer(c *battleship_client.BattleshipHTTPClient) error {
	playerList, err := c.GetPlayersList(WaitingPlayersEndpoint)
	if err != nil {
		return err
	}

	for _, data := range *playerList {
		fmt.Println(data.GameStatus, " | ", data.Nick)
	}
	return nil
}

func StartNewGame(c *battleship_client.BattleshipHTTPClient) error {
	app := New(c)
	err := app.Run()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
