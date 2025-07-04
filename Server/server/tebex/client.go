package tebex

import (
	"bytes"
	"log/slog"
	"net/http"
	"time"

	"gopkg.in/square/go-jose.v2/json"
)

// baseURL is the base URL for the Tebex API.
const baseURL = "https://plugin.tebex.io/"

// Client is a client used to access the Tebex API. It can be used to get information about orders, players, bans, gift
// cards, and more.
type Client struct {
	secret string
	period time.Duration

	log *slog.Logger
	c   *http.Client

	s chan struct{}
}

// NewClient creates a new Client with the given secret. The *http.Client will be defaulted to http.DefaultClient.
// The period provided will be used to determine how often the client will check for new offline commands.
func NewClient(log *slog.Logger, period time.Duration, secret string) *Client {
	c := &Client{
		s: make(chan struct{}),
		c: http.DefaultClient,

		period: period,
		secret: secret,
		log:    log,
	}
	go c.startTicking()
	return c
}

// Information returns the name and domain of the Tebex account associated with the client.
func (c *Client) Information() (string, string, error) {
	var result struct {
		Account struct {
			Name   string `json:"name"`
			Domain string `json:"domain"`
		} `json:"account"`
	}
	err := c.get("information", &result)
	if err != nil {
		return "", "", err
	}
	return result.Account.Name, result.Account.Domain, nil
}

// Close closes the client.
func (c *Client) Close() error {
	close(c.s)
	return nil
}

// startTicking starts a ticker that will check for new offline commands every thirty seconds.
func (c *Client) startTicking() {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}

	ticker := time.NewTicker(c.period)
	defer ticker.Stop()

	for {
		select {
		case <-ch:
			c.ExecuteOfflineCommands()
		case <-ticker.C:
			ch <- struct{}{}
		case <-c.s:
			return
		}
	}
}

// get performs a GET request using the *http.Client attached. It will always include the X-Tebex-Secret header.
func (c *Client) get(endpoint string, o any) error {
	req, err := http.NewRequest(http.MethodGet, baseURL+endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Tebex-Secret", c.secret)
	req.Header.Set("User-Agent", "Tebex for Dragonfly")
	res, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(o)
}

// delete performs a DELETE request using the *http.Client attached. It will always include the X-Tebex-Secret header.
func (c *Client) delete(endpoint string, params map[string]any) error {
	b, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, baseURL+endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("X-Tebex-Secret", c.secret)
	req.Header.Set("User-Agent", "Tebex for Dragonfly")
	req.Header.Set("Content-Type", "application/json")
	res, err := c.c.Do(req)
	if err != nil {
		return err
	}
	return res.Body.Close()
}
