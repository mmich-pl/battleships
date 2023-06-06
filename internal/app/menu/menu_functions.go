package menu

import (
	"battleships/internal/battleship_client"
	. "battleships/internal/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
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

func PrintPlayerStatistics(c battleship_client.BattleshipClient) error {
	input, _ := GetPlayerInput("get the stats for player:")

	s, err := c.GetPlayerStatistic(strings.TrimSpace(input))
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

func PrintTopTenPlayerStatistics(c battleship_client.BattleshipClient) error {
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', tabwriter.AlignRight)
	stats, err := c.GetStatistic()
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

func ListPlayer(c battleship_client.BattleshipClient) error {
	playerList, err := c.GetPlayersList()
	if err != nil {
		return err
	}

	for _, data := range *playerList {
		fmt.Println(data.GameStatus, " | ", data.Nick)
	}
	return nil
}
