package battlehip_client

import (
	"battleships/internal/models"
	"battleships/pkg/base_client"
	"fmt"
	"log"
	"net/http"
	"time"
)

type BattleshipClient interface {
	InitGame(endpoint, nick, desc, targetNick, string, wpbot bool) error
	GameStatus()
}

type BattleshipHTTPClient struct {
	client *base_client.BaseHTTPClient
	token  string
}

func NewBattleshipClient(baseURL string, responseTimeout, connectionTimeout time.Duration) *BattleshipHTTPClient {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "go")
	client := base_client.NewBuilder().
		SetHeaders(headers).
		SetConnectionTimeout(connectionTimeout).
		SetResponseTime(responseTimeout).
		SetBaseURL(baseURL).
		Build()
	return &BattleshipHTTPClient{client: &client}
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
