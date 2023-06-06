package battleship_client

import (
	"battleships/internal/models"
	"battleships/pkg/base_client"
	"fmt"
	"net/http"
	"time"
)

const (
	InitEndpoint           = "/game"
	GameStatusEndpoint     = "/game"
	BoardEndpoint          = "/game/board"
	OpponentDescription    = "/game/desc"
	FireEndpoint           = "/game/fire"
	WaitingPlayersEndpoint = "/lobby"
	RefreshEndpoint        = "/game/refresh"
	StatsEndpoint          = "/stats"
	AbandonEndpoint        = "/game/abandon"
)

var (
	defaultRetryWaitMin = 0
	defaultRetryWaitMax = 10
	defaultRetryMax     = 4
)

//go:generate mockery --name BattleshipClient
type BattleshipClient interface {
	InitGame(nick, desc, targetNick string, coords []string, wpbot bool) error
	Description() (*models.DescriptionResponse, error)
	GameStatus() (*models.StatusResponse, error)
	Board() ([]string, error)
	Fire(coords string) (*models.ShootResult, error)
	GetPlayersList() (*[]models.WaitingPlayerData, error)
	RefreshSession() error
	GetStatistic() (*models.StatsResponse, error)
	GetPlayerStatistic(nick string) (*models.PlayerStatsResponse, error)
	AbandonGame() error
	GetToken() string
	ResetToken()
}

type BattleshipHTTPClient struct {
	client *base_client.BaseHTTPClient
	token  string
}

func NewBattleshipClient(baseURL string, responseTimeout, connectionTimeout time.Duration) *BattleshipHTTPClient {
	headers := http.Header{
		"Content-Type": {"application/json"},
	}

	clientConfig := base_client.HTTPClientConfig{
		BaseUrl:           baseURL,
		ConnectionTimeout: connectionTimeout * time.Second,
		ResponseTimeout:   responseTimeout * time.Second,
		Headers:           headers,
		//ProxyAddress:      "http://127.0.0.1:8900",
		RetryWaitMin:  defaultRetryWaitMin,
		RetryWaitMax:  defaultRetryWaitMax,
		RetryMax:      defaultRetryMax,
		CheckForRetry: base_client.DefaultRetryPolicy,
		Backoff:       base_client.DefaultBackoff,
	}

	return &BattleshipHTTPClient{client: base_client.New(clientConfig)}
}

func (b *BattleshipHTTPClient) InitGame(nick, desc, targetNick string, coords []string, wpbot bool) error {
	payload := models.InitialPayload{
		Coords:     coords,
		Desc:       desc,
		Nick:       nick,
		TargetNick: targetNick,
		Wpbot:      wpbot,
	}

	resp, err := b.client.Post(InitEndpoint, payload, b.client.Config.Headers)
	if err != nil {
		return fmt.Errorf("failed to perform POST rquest: %w", err)
	}

	b.token = resp.Headers.Get("X-Auth-Token")
	b.client.AddHeader("X-Auth-token", b.token)
	return nil
}

func (b *BattleshipHTTPClient) Description() (*models.DescriptionResponse, error) {
	resp, err := b.client.Get(OpponentDescription, b.client.Config.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET request: %w", err)
	}
	var descriptions models.DescriptionResponse
	if err = resp.UnmarshalJson(&descriptions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return &descriptions, nil
}

func (b *BattleshipHTTPClient) GameStatus() (*models.StatusResponse, error) {
	resp, err := b.client.Get(GameStatusEndpoint, b.client.Config.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET request: %w", err)
	}
	var status models.StatusResponse
	if err = resp.UnmarshalJson(&status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return &status, nil
}

func (b *BattleshipHTTPClient) Board() ([]string, error) {
	type Board struct {
		Board []string `json:"board"`
	}
	resp, err := b.client.Get(BoardEndpoint, b.client.Config.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET request: %w", err)
	}
	board := Board{Board: make([]string, 20)}
	if err := resp.UnmarshalJson(&board); err != nil {
		return nil, fmt.Errorf("failed to unmarshal board: %w", err)
	}
	return board.Board, nil
}

func (b *BattleshipHTTPClient) Fire(coords string) (*models.ShootResult, error) {
	payload := models.Shoot{Coord: coords}
	resp, err := b.client.Post(FireEndpoint, payload, b.client.Config.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform POST rquest: %w", err)
	}

	var result models.ShootResult
	if err = resp.UnmarshalJson(&result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hit result: %w", err)
	}
	return &result, nil
}

func (b *BattleshipHTTPClient) GetPlayersList() (*[]models.WaitingPlayerData, error) {
	resp, err := b.client.Get(WaitingPlayersEndpoint, b.client.Config.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get response: %w", err)
	}

	var result []models.WaitingPlayerData
	if err = resp.UnmarshalJson(&result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal waiting players request: %w", err)
	}
	return &result, nil
}

func (b *BattleshipHTTPClient) RefreshSession() error {
	resp, err := b.client.Get(RefreshEndpoint, b.client.Config.Headers)
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (b *BattleshipHTTPClient) GetStatistic() (*models.StatsResponse, error) {
	resp, err := b.client.Get(StatsEndpoint, b.client.Config.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get response: %w", err)
	}
	var result models.StatsResponse
	if err = resp.UnmarshalJson(&result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal waiting players request: %w", err)
	}
	return &result, nil
}

func (b *BattleshipHTTPClient) GetPlayerStatistic(nick string) (*models.PlayerStatsResponse, error) {
	resp, err := b.client.Get(StatsEndpoint+"/"+nick, b.client.Config.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get response: %w", err)
	}
	var result models.PlayerStatsResponse
	if err = resp.UnmarshalJson(&result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal waiting players request: %w", err)
	}
	return &result, nil
}

func (b *BattleshipHTTPClient) AbandonGame() error {
	resp, err := b.client.Delete(AbandonEndpoint, b.client.Config.Headers)
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (b *BattleshipHTTPClient) GetToken() string {
	return b.token
}

func (b *BattleshipHTTPClient) ResetToken() {
	b.token = ""
}
