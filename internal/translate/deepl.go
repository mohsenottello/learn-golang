package translate

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrDisabled = errors.New("translation is not configured")

type Client struct {
	key  string
	base string
	http *http.Client
}

func New(key string) *Client {
	base := "https://api.deepl.com"
	if strings.HasSuffix(key, ":fx") {
		base = "https://api-free.deepl.com"
	}
	return &Client{key: key, base: base, http: &http.Client{Timeout: 15 * time.Second}}
}

func (c *Client) Enabled() bool { return c != nil && c.key != "" }

// Translate sends several texts in one call and returns them in the same order.
func (c *Client) Translate(ctx context.Context, texts []string, source, target string) ([]string, error) {
	if !c.Enabled() {
		return nil, ErrDisabled
	}
	form := url.Values{"source_lang": {strings.ToUpper(source)}, "target_lang": {strings.ToUpper(target)}}
	for _, t := range texts {
		form.Add("text", t)
	}
	req, _ := http.NewRequestWithContext(ctx, "POST", c.base+"/v2/translate", strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "DeepL-Auth-Key "+c.key)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 456 {
		return nil, errors.New("DeepL quota exhausted for this month")
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("DeepL error: " + resp.Status)
	}
	var out struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	res := make([]string, len(out.Translations))
	for i, t := range out.Translations {
		res[i] = t.Text
	}
	return res, nil
}
