package models

type InitialPayload struct {
	Coords     []string `json:"coords,omitempty"`
	Desc       string   `json:"desc"`
	Nick       string   `json:"nick"`
	TargetNick string   `json:"target_nick,omitempty"`
	Wpbot      bool     `json:"wpbot"`
}

type StatusResponse struct {
	Desc                string   `json:"desc"`
	GameStatus          string   `json:"game_status"`
	LastGameStatus      string   `json:"last_game_status"`
	Nick                string   `json:"nick"`
	OpponentDescription string   `json:"opp_desc"`
	OpponentShots       []string `json:"opp_shots"`
	Opponent            string   `json:"opponent"`
	ShouldFire          bool     `json:"should_fire"`
	Timer               int32    `json:"timer"`
}
