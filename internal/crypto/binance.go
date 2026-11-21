// Package crypto fetches spot prices from Binance's public API (no key needed).
package crypto

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct{ http *http.Client }

func New() *Client { return &Client{http: &http.Client{Timeout: 10 * time.Second}} }

// Prices returns the USDT price of each symbol ("BTC" → 94210.5). Unknown symbols are omitted.
func (c *Client) Prices(ctx context.Context, symbols []string) (map[string]float64, error) {
	out := map[string]float64{}
	if len(symbols) == 0 {
		return out, nil
	}
	pairs := make([]string, len(symbols))
	for i, s := range symbols {
		pairs[i] = strings.ToUpper(s) + "USDT"
	}
	list, _ := json.Marshal(pairs) // ["BTCUSDT","ETHUSDT"]

	req, _ := http.NewRequestWithContext(ctx, "GET",
		"https://api.binance.com/api/v3/ticker/price?symbols="+url.QueryEscape(string(list)), nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// Binance answers 400 if *any* symbol is unknown; fall back to one call per symbol.
		if resp.StatusCode == http.StatusBadRequest {
			return c.pricesOneByOne(ctx, symbols)
		}
		return nil, fmt.Errorf("binance: %s", resp.Status)
	}
	var rows []struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, err
	}
	for _, r := range rows {
		if p, err := strconv.ParseFloat(r.Price, 64); err == nil {
			out[strings.TrimSuffix(r.Symbol, "USDT")] = p
		}
	}
	return out, nil
}

// Price returns one symbol's price; used to validate a symbol when a watch is created.
func (c *Client) Price(ctx context.Context, symbol string) (float64, error) {
	m, err := c.pricesOneByOne(ctx, []string{symbol})
	if err != nil {
		return 0, err
	}
	p, ok := m[strings.ToUpper(symbol)]
	if !ok {
		return 0, fmt.Errorf("unknown symbol %s", symbol)
	}
	return p, nil
}

func (c *Client) pricesOneByOne(ctx context.Context, symbols []string) (map[string]float64, error) {
	out := map[string]float64{}
	for _, s := range symbols {
		s = strings.ToUpper(s)
		req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.binance.com/api/v3/ticker/price?symbol="+s+"USDT", nil)
		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		var r struct {
			Price string `json:"price"`
		}
		if resp.StatusCode == http.StatusOK {
			if json.NewDecoder(resp.Body).Decode(&r) == nil {
				if p, err := strconv.ParseFloat(r.Price, 64); err == nil {
					out[s] = p
				}
			}
		}
		resp.Body.Close() // unknown symbol → simply not in the map
	}
	return out, nil
}
