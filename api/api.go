//go:generate go run ../cmd/apigen.go
package api

import (
	"context"
	"errors"
)

const (
	GetPriceEndpoint     = "/price/{assetBase}/{assetQuote}"
	GetExchangesEndpoint = "/exchanges"
	GetAssetsEndpoint    = "/assets"
)

type API interface {
	GetPrice(ctx context.Context, req GetPriceRequest) (*GetPriceResponse, error)
	//GetExchanges(ctx context.Context, req GetExchangesRequest) (*GetExchangesResponse, error)
}

type api struct {
}

func NewAPI() API {
	return &api{}
}

const BaseURL = "https://rest.coinapi.io/v1"

type GetPriceRequest struct {
	AssetBase  string `json:"assetBase"`
	AssetQuote string `json:"assetQuote"`
}

func (r GetPriceRequest) Validate() error {
	if r.AssetBase == "" {
		return errors.New("assetBase cannot be null")
	}
	if r.AssetQuote == "" {
		return errors.New("assetQuote cannot be null")
	}
	return nil
}

type GetPriceResponse struct {
	Rate  float64 `json:"rate,omitempty"`
	Error string  `json:"error,omitempty"`
}

func (a *api) GetPrice(ctx context.Context, req GetPriceRequest) (*GetPriceResponse, error) {
	//	client := &http.Client{}
	//	endpoint := fmt.Sprintf("%s/exchangerate/%s/%s", BaseURL, req.AssetBase, req.AssetQuote)
	//	httpReq, err := http.NewRequest(http.MethodGet, endpoint, nil)
	//	if err != nil {
	//		return nil, err
	//	}
	//	httpReq.Header.Set("X-CoinAPI-Key", "C93B0915-91FF-4CD3-B07F-BA32376F900F")
	//	r, _ := client.Do(httpReq)
	//
	//	var coinResp coinapi.ExchangeRate
	//	dec := json.NewDecoder(r.Body)
	//	err = dec.Decode(&coinResp)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return &GetPriceResponse{
	//		Rate:  coinResp.Rate,
	//		Error: coinResp.Error,
	//	}, nil
	return &GetPriceResponse{
		Rate: 1.0,
	}, nil
}

// type GetExchangesRequest struct {
// }
//
// func (r GetExchangesRequest) Validate() error {
// 	return nil
// }
//
// type GetExchangesResponse struct {
// }
//
// func (a *api) GetExchanges(ctx context.Context, req GetExchangesRequest) (*GetExchangesResponse, error) {
// 	return nil, nil
// }
