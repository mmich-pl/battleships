package battlehip_client

import (
	"battleships/battlehip_client/models"
	"fmt"
	"log"
	"net/http"
	"time"
)

type BattleshipHTTPClient interface {
	InitGame(endpoint, nick, desc, targetNick, string, wpbot bool) error
}

type BattleshipClient struct {
	client *Client
	token  string
}

func NewBattleshipClient(baseURL string, responseTimeout, connectionTimeout time.Duration) *BattleshipClient {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "go")
	client := NewBuilder().
		SetHeaders(headers).
		SetConnectionTimeout(connectionTimeout).
		SetResponseTime(responseTimeout).
		SetBaseURL(baseURL).
		Build()
	return &BattleshipClient{client: &client}
}

func (b *BattleshipClient) InitGame(endpoint, nick, desc, targetNick string, wpbot bool) error {
	payload := models.InitialPayload{
		Coords:     nil,
		Desc:       nick,
		Nick:       nick,
		TargetNick: targetNick,
		Wpbot:      wpbot,
	}

	resp, err := b.client.Post(endpoint, payload, b.client.builder.headers)
	if err != nil {
		return fmt.Errorf("failed to perform POST rquest: %w", err)
	}

	b.token = resp.Headers.Get("X-Auth-Token")
	log.Printf("Client's token: %s", b.token)
	return nil
}
