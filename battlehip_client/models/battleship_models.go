package models

type InitialPayload struct {
	Coords     []string `json:"coords,omitempty"`
	Desc       string   `json:"desc"`
	Nick       string   `json:"nick"`
	TargetNick string   `json:"target_nick,omitempty"`
	Wpbot      bool     `json:"wpbot"`
}
