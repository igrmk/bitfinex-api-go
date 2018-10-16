package rest

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	bitfinex "github.com/bitfinexcom/bitfinex-api-go/v2"
)

type TickerService struct {
	Synchronous
}

func (t *TickerService) All() (*bitfinex.TickerSnapshot, error) {
	req := NewRequestWithMethod("tickers", "GET")
	req.Params = make(url.Values)
	req.Params.Add("symbols", "ALL")
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
		if len(arr) != 11 && len(arr) != 12 {
			return nil, errors.New("invalid length of ticker")
		}
		symbol, ok := arr[0].(string)
		symbol = strings.ToLower(symbol[1:])
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
		var entry *bitfinex.Ticker
		switch len(arr) {
		case 11:
			entry = &bitfinex.Ticker{
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
		case 12:
			entry = &bitfinex.Ticker{
				Symbol:          symbol,
				FRR:             sub[0],
				Bid:             sub[1],
				BidSize:         sub[2],
				Ask:             sub[3],
				AskSize:         sub[4],
				DailyChange:     sub[5],
				DailyChangePerc: sub[6],
				LastPrice:       sub[7],
				Volume:          sub[8],
				High:            sub[9],
				Low:             sub[10],
			}
		}
		tickers[i] = entry
	}
	return &bitfinex.TickerSnapshot{Snapshot: tickers}, nil
}
