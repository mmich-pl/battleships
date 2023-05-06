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
	GameStatus     string   `json:"game_status,omitempty"`
	LastGameStatus string   `json:"last_game_status,omitempty"`
	OpponentShots  []string `json:"opp_shots,omitempty"`
	ShouldFire     bool     `json:"should_fire,omitempty"`
	Timer          int32    `json:"timer,omitempty"`
}

type Shoot struct {
	Coordinate string `json:"coord"`
}

type ShootResult struct {
	Result string `json:"result"`
}
