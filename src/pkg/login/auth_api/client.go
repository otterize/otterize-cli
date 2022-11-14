package auth_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	address string
}

type Auth0Config struct {
	Domain   string `json:"domain"`
	Audience string `json:"audience"`
	ClientId string `json:"client_id"`
}

func NewClient(address string) *Client {
	address = address + "/auth"
	return &Client{
		address: address,
	}
}

func (c *Client) GetAuth0Config(ctx context.Context) (Auth0Config, error) {
	confUrl := fmt.Sprintf("%s/auth0-config", c.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, confUrl, nil)
	if err != nil {
		return Auth0Config{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Auth0Config{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Auth0Config{}, fmt.Errorf("HTTP status code %d", resp.StatusCode)
	}

	var r Auth0Config
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Auth0Config{}, err
	}

	return r, nil
}
