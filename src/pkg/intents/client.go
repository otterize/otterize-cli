package intents

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
)

type Client struct {
	c *cloudclient.Client
}

func NewClientFromToken(address string, token string) *Client {
	cloud := cloudclient.NewClientFromToken(address, token)
	return &Client{c: cloud}
}

func (c *Client) ReportDiscoveredIntents(ctx context.Context, envId string, source string, intents []IntentInput) error {
	_, err := reportDiscoveredIntents(ctx, c.c.Client, envId, source, intents)
	return err
}
