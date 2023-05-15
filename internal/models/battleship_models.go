package models

type InitialPayload struct {
	Coords     []string `json:"coords,omitempty"`
	Desc       string   `json:"desc"`
	Nick       string   `json:"nick"`
	TargetNick string   `json:"target_nick,omitempty"`
	Wpbot      bool     `json:"wpbot"`
}

type DescriptionResponse struct {
	Desc                string `json:"desc,omitempty"`
	Nick                string `json:"nick,omitempty"`
	OpponentDescription string `json:"opp_desc,omitempty"`
	Opponent            string `json:"opponent,omitempty"`
}

type StatusResponse struct {
	GameStatus     string   `json:"game_status"`
	LastGameStatus string   `json:"last_game_status"`
	Nick           string   `json:"nick"`
	OpponentShots  []string `json:"opp_shots"`
	Opponent       string   `json:"opponent"`
	ShouldFire     bool     `json:"should_fire"`
	Timer          int      `json:"timer"`
}

type Shoot struct {
	Coord string `json:"coord"`
}

type ShootResult struct {
	Result string `json:"result"`
}

type WaitingPlayerData struct {
	GameStatus string `json:"game_status"`
	Nick       string `json:"nick"`
}

type PlayerStatistics struct {
	Stats struct {
		Games  int    `json:"games"`
		Nick   string `json:"nick"`
		Points int    `json:"points"`
		Rank   int    `json:"rank"`
		Wins   int    `json:"wins"`
	} `json:"stats"`
}
