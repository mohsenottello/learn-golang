package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	token string
	http  *http.Client
}

func New(token string) *Client {
	return &Client{token: token, http: &http.Client{Timeout: 40 * time.Second}}
}

func (c *Client) Enabled() bool { return c != nil && c.token != "" }

func (c *Client) call(ctx context.Context, method string, params url.Values, out any) error {
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.telegram.org/bot"+c.token+"/"+method, nil)
	req.URL.RawQuery = params.Encode()
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var env struct {
		OK          bool            `json:"ok"`
		Description string          `json:"description"`
		Result      json.RawMessage `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return err
	}
	if !env.OK {
		return errors.New("telegram: " + env.Description)
	}
	if out != nil {
		return json.Unmarshal(env.Result, out)
	}
	return nil
}

func (c *Client) Send(ctx context.Context, chatID int64, html string) error {
	return c.call(ctx, "sendMessage", url.Values{
		"chat_id": {strconv.FormatInt(chatID, 10)}, "text": {html}, "parse_mode": {"HTML"},
	}, nil)
}

type Update struct {
	ID      int64 `json:"update_id"`
	Message *struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

// Updates long-polls for new messages (blocks up to 30 s).
func (c *Client) Updates(ctx context.Context, offset int64) ([]Update, error) {
	var out []Update
	err := c.call(ctx, "getUpdates", url.Values{
		"timeout": {"30"}, "offset": {strconv.FormatInt(offset, 10)}, "allowed_updates": {`["message"]`},
	}, &out)
	return out, err
}
