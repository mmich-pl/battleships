package battlehip_client

import (
	"battleships/internal/models"
	"battleships/pkg/base_client"
	"fmt"
	"time"
)

//go:generate mockery --name BattleshipClient
type BattleshipClient interface {
	InitGame(endpoint, nick, desc, targetNick string, wpbot bool) error
	Description(endpoint string) (*models.DescriptionResponse, error)
	GameStatus(endpoint string) (*models.StatusResponse, error)
	Board(endpoint string) ([]string, error)
	Fire(endpoint, coords string) (*models.ShootResult, error)
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
		SetResponseTimeout(responseTimeout).
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
	return nil
}

func (b *BattleshipHTTPClient) Description(endpoint string) (*models.DescriptionResponse, error) {
	resp, err := b.client.Get(endpoint, b.client.Builder.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET request: %w", err)
	}
	var descriptions models.DescriptionResponse
	if err := resp.UnmarshalJson(&descriptions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return &descriptions, nil
}

func (b *BattleshipHTTPClient) GameStatus(endpoint string) (*models.StatusResponse, error) {
	resp, err := b.client.Get(endpoint, b.client.Builder.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET request: %w", err)
	}
	var status models.StatusResponse
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
	return board.Board, nil
}

func (b *BattleshipHTTPClient) Fire(endpoint, coords string) (*models.ShootResult, error) {
	payload := models.Shoot{Coord: coords}
	resp, err := b.client.Post(endpoint, payload, b.client.Builder.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform POST rquest: %w", err)
	}

	var result models.ShootResult
	if err = resp.UnmarshalJson(&result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hit result: %w", err)
	}
	return &result, nil
}
