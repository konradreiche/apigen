package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/konradreiche/apigen/coinapi"
)

type API interface {
	GetPrice(ctx context.Context, req GetPriceRequest) (*GetPriceResponse, error)
}

type api struct {
}

func NewAPI() API {
	return &api{}
}

const BaseURL = "https://rest.coinapi.io/v1"

type Error struct {
	Error string `json:"error,omitempty"`
}

type GetPriceRequest struct {
	AssetBase  string `json:"assetBase"`
	AssetQuote string `json:"assetQuote"`
}

type GetPriceResponse struct {
	Rate float64 `json:"rate,omitempty"`
	Error
}

func (a *api) GetPrice(ctx context.Context, req GetPriceRequest) (*GetPriceResponse, error) {
	client := &http.Client{}
	endpoint := fmt.Sprintf("%s/exchangerate/%s/%s", BaseURL, req.AssetBase, req.AssetQuote)
	httpReq, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("X-CoinAPI-Key", "C93B0915-91FF-4CD3-B07F-BA32376F900F")
	r, _ := client.Do(httpReq)

	var coinResp coinapi.ExchangeRate
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&coinResp)
	if err != nil {
		return nil, err
	}
	return &GetPriceResponse{
		Rate:  coinResp.Rate,
		Error: Error{Error: coinResp.Error},
	}, nil
}
