package rest

import (
	"errors"
	"fmt"

	bitfinex "github.com/igrmk/bitfinex-api-go/v2"
)

type TickerService struct {
	Synchronous
}

func (t *TickerService) All() (*bitfinex.TickerSnapshot, error) {
	req := NewRequestWithMethod("tickers", "GET")
	raw, err := t.Request(req)

	if err != nil {
		return nil, err
	}

	tickers := make([]*bitfinex.Ticker, len(raw))
	for i, ifacearr := range raw {
		arr, ok := ifacearr.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expecting array, got %T", ifacearr)
		}
		if len(arr) != 11 {
			return nil, errors.New("invalid length of ticker")
		}
		symbol, ok := arr[0].(string)
		if !ok {
			return nil, fmt.Errorf("expecting string, got %T", arr[0])
		}
		sub := make([]float64, len(arr)-1)
		for j, iface := range arr[1:] {
			flt, ok := iface.(float64)
			if !ok {
				return nil, fmt.Errorf("expecting float64, got %T", iface)
			}
			sub[j] = flt
		}
		entry := &bitfinex.Ticker{
			Symbol:          symbol,
			Bid:             sub[0],
			BidSize:         sub[1],
			Ask:             sub[2],
			AskSize:         sub[3],
			DailyChange:     sub[4],
			DailyChangePerc: sub[5],
			LastPrice:       sub[6],
			Volume:          sub[7],
			High:            sub[8],
			Low:             sub[9],
		}
		tickers[i] = entry
	}
	return &bitfinex.TickerSnapshot{Snapshot: tickers}, nil
}
