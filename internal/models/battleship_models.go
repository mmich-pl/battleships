package models

type InitialPayload struct {
	Coords     []string `json:"coords,omitempty"`
	Desc       string   `json:"desc"`
	Nick       string   `json:"nick"`
	TargetNick string   `json:"target_nick,omitempty"`
	Wpbot      bool     `json:"wpbot"`
}

type FullStatusResponse struct {
	Desc                string   `json:"desc,omitempty"`
	GameStatus          string   `json:"game_status,omitempty"`
	LastGameStatus      string   `json:"last_game_status,omitempty"`
	Nick                string   `json:"nick,omitempty"`
	OpponentDescription string   `json:"opp_desc,omitempty"`
	OpponentShots       []string `json:"opp_shots,omitempty"`
	Opponent            string   `json:"opponent,omitempty"`
	ShouldFire          bool     `json:"should_fire,omitempty"`
	Timer               int32    `json:"timer,omitempty"`
}

type PartialStatusResponse struct {
	GameStatus     string   `json:"game_status,omitempty"`
	LastGameStatus string   `json:"last_game_status,omitempty"`
	OpponentShots  []string `json:"opp_shots,omitempty"`
	ShouldFire     bool     `json:"should_fire,omitempty"`
	Timer          int32    `json:"timer,omitempty"`
}
