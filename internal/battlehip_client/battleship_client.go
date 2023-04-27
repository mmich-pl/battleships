package battlehip_client

import (
	"battleships/internal/models"
	"battleships/pkg/base_client"
	"fmt"
	"log"
	"time"
)

//go:generate mockery --name BattleshipClient
type BattleshipClient interface {
	InitGame(endpoint, nick, desc, targetNick string, wpbot bool) error
	FullGameStatus(endpoint string) (*models.FullStatusResponse, error)
	Board(endpoint string) ([]string, error)
}

type BattleshipHTTPClient struct {
	client *base_client.BaseHTTPClient
	token  string
}

func NewBattleshipClient(baseURL string, responseTimeout, connectionTimeout time.Duration) *BattleshipHTTPClient {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	client := base_client.NewBuilder().
		SetHeaderFromMap(headers).
		SetConnectionTimeout(connectionTimeout).
		SetResponseTime(responseTimeout).
		SetBaseURL(baseURL).
		Build()
	return &BattleshipHTTPClient{client: client}
}

func (b *BattleshipHTTPClient) InitGame(endpoint, nick, desc, targetNick string, wpbot bool) error {
	payload := models.InitialPayload{
		Coords:     nil,
		Desc:       desc,
		Nick:       nick,
		TargetNick: targetNick,
		Wpbot:      wpbot,
	}

	resp, err := b.client.Post(endpoint, payload, b.client.Builder.Headers)
	if err != nil {
		return fmt.Errorf("failed to perform POST rquest: %w", err)
	}

	b.token = resp.Headers.Get("X-Auth-Token")
	b.client.Builder.AddHeader("X-Auth-Token", b.token)
	log.Printf("BaseHTTPClient's token: %s", b.token)
	return nil
}

func (b *BattleshipHTTPClient) FullGameStatus(endpoint string) (*models.FullStatusResponse, error) {
	resp, err := b.client.Get(endpoint, b.client.Builder.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET request: %w", err)
	}
	var status models.FullStatusResponse
	if err := resp.UnmarshalJson(&status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return &status, nil
}

func (b *BattleshipHTTPClient) Board(endpoint string) ([]string, error) {
	type Board struct {
		Board []string `json:"board"`
	}
	resp, err := b.client.Get(endpoint, b.client.Builder.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET request: %w", err)
	}
	board := Board{Board: make([]string, 20)}
	if err := resp.UnmarshalJson(&board); err != nil {
		return nil, fmt.Errorf("failed to unmarshal board: %w", err)
	}
	log.Println(board)
	return board.Board, nil
}

func (b *BattleshipHTTPClient) Fire(endpoint, coords string) (string, error) {
	resp, err := b.client.Post(endpoint, coords, b.client.Builder.Headers)
	if err != nil {
		return "", fmt.Errorf("failed to perform POST rquest: %w", err)
	}

	result := struct {
		Result string `json:"result"`
	}{}
	if err := resp.UnmarshalJson(&resp); err != nil {
		return "", fmt.Errorf("failed to unmarshal hit result: %w", err)
	}
	log.Print(result.Result)
	return result.Result, nil
}
