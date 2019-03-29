package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/konradreiche/apigen/api"
)

func (c *Client) GetPrice(ctx context.Context, req api.GetPriceRequest) (*api.GetPriceResponse, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/price"), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.conn.Do(r)
	if err != nil {
		return nil, err
	}
	var result api.GetPriceResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, err
}
