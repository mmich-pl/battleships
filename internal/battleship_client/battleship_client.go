package battleship_client

import (
	"battleships/internal/models"
	"battleships/pkg/base_client"
	"fmt"
	"net/http"
	"time"
)

//go:generate mockery --name BattleshipClient
type BattleshipClient interface {
	InitGame(endpoint, nick, desc, targetNick string, coords []string, wpbot bool) error
	Description(endpoint string) (*models.DescriptionResponse, error)
	GameStatus(endpoint string) (*models.StatusResponse, error)
	Board(endpoint string) ([]string, error)
	Fire(endpoint, coords string) (*models.ShootResult, error)
	GetPlayersList(endpoint string) (*[]models.WaitingPlayerData, error)
	RefreshSession(endpoint string) error
	GetStatistic(endpoint string) (*models.StatsResponse, error)
	GetPlayerStatistic(endpoint, nick string) (*models.PlayerStatsResponse, error)
	AbandonGame(endpoint string) error
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

func (b *BattleshipHTTPClient) InitGame(endpoint, nick, desc, targetNick string, coords []string, wpbot bool) error {
	payload := models.InitialPayload{
		Coords:     coords,
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
	if err = resp.UnmarshalJson(&descriptions); err != nil {
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

func (b *BattleshipHTTPClient) GetPlayersList(endpoint string) (*[]models.WaitingPlayerData, error) {
	resp, err := b.client.Get(endpoint, b.client.Builder.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get response: %w", err)
	}

	var result []models.WaitingPlayerData
	if err = resp.UnmarshalJson(&result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal waiting players request: %w", err)
	}
	return &result, nil
}

func (b *BattleshipHTTPClient) RefreshSession(endpoint string) error {
	resp, err := b.client.Get(endpoint, b.client.Builder.Headers)
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (b *BattleshipHTTPClient) GetStatistic(endpoint string) (*models.StatsResponse, error) {
	resp, err := b.client.Get(endpoint, b.client.Builder.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get response: %w", err)
	}
	var result models.StatsResponse
	if err = resp.UnmarshalJson(&result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal waiting players request: %w", err)
	}
	return &result, nil
}

func (b *BattleshipHTTPClient) GetPlayerStatistic(endpoint, nick string) (*models.PlayerStatsResponse, error) {
	resp, err := b.client.Get(endpoint+"/"+nick, b.client.Builder.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get response: %w", err)
	}
	var result models.PlayerStatsResponse
	if err = resp.UnmarshalJson(&result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal waiting players request: %w", err)
	}
	return &result, nil
}

func (b *BattleshipHTTPClient) AbandonGame(endpoint string) error {
	resp, err := b.client.Delete(endpoint, b.client.Builder.Headers)
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
