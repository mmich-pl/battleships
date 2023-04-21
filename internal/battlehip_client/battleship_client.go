package battlehip_client

import (
	"battleships/internal/models"
	"battleships/pkg/base_client"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type BattleshipClient interface {
	InitGame(endpoint, nick, desc, targetNick, string, wpbot bool) error
	GameStatus(endpoint string) (*models.StatusData, error)
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
	log.Printf("BaseHTTPClient's token: %s", b.token)
	return nil
}

func (b *BattleshipHTTPClient) GameStatus(endpoint string) (*models.StatusData, error) {
	b.client.Builder.AddHeader("X-Auth-Token", b.token)
	resp, err := b.client.Get(endpoint, b.client.Builder.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET request: %w", err)
	}
	var status *models.StatusData
	if err := json.Unmarshal(resp.Body, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return status, nil
}
