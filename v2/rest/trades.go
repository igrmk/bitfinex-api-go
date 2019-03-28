package rest

import (
	"encoding/json"
	"path"

	"github.com/igrmk/bitfinex-api-go/v2"
)

// TradeService manages the Trade endpoint.
type TradeService struct {
	requestFactory
	Synchronous
}

// All returns all orders for the authenticated account.
func (s *TradeService) All(symbol string) (*bitfinex.TradeSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(path.Join("trades", symbol, "hist"), map[string]interface{}{"start": nil, "end": nil, "limit": nil})
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)

	if err != nil {
		return nil, err
	}

	dat := make([][]json.Number, 0)
	for _, r := range raw {
		if f, ok := r.([]json.Number); ok {
			dat = append(dat, f)
		}
	}

	os, err := bitfinex.NewTradeSnapshotFromRaw(symbol, dat)
	if err != nil {
		return nil, err
	}
	return os, nil
}
