package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	bitfinex "github.com/igrmk/bitfinex-api-go/v2"
)

type TickerService struct {
	Synchronous
}

var decValOrZero = bitfinex.DecValOrZero
var int64ValOrZero = bitfinex.Int64ValOrZero

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
		symbol, ok := arr[0].(string)
		if !ok {
			return nil, fmt.Errorf("expecting string, got %T", arr[0])
		}
		if len(symbol) <= 1 || (symbol[0] != 't' && symbol[0] != 'f') {
			return nil, errors.New("invalid symbol")
		}
		if (symbol[0] == 't' && len(arr) < 11) || (symbol[0] == 'f' && len(arr) < 14) {
			return nil, errors.New("invalid length of ticker")
		}
		sub := make([]json.Number, len(arr)-1)
		for j, iface := range arr[1:] {
			if iface == nil {
				sub[j] = "0"
				continue
			}
			flt, ok := iface.(json.Number)
			if !ok {
				return nil, fmt.Errorf("expecting a number, got %T", iface)
			}
			sub[j] = flt
		}
		var entry *bitfinex.Ticker
		switch symbol[0] {
		case 't':
			entry = &bitfinex.Ticker{
				Symbol:          strings.ToLower(symbol[1:]),
				Bid:             decValOrZero(sub[0]),
				BidSize:         decValOrZero(sub[1]),
				Ask:             decValOrZero(sub[2]),
				AskSize:         decValOrZero(sub[3]),
				DailyChange:     decValOrZero(sub[4]),
				DailyChangePerc: decValOrZero(sub[5]),
				LastPrice:       decValOrZero(sub[6]),
				Volume:          decValOrZero(sub[7]),
				High:            decValOrZero(sub[8]),
				Low:             decValOrZero(sub[9]),
			}
		case 'f':
			entry = &bitfinex.Ticker{
				Symbol:          strings.ToLower(symbol[1:]),
				FRR:             decValOrZero(sub[0]),
				Bid:             decValOrZero(sub[1]),
				BidSize:         decValOrZero(sub[2]),
				BidPeriod:       int64ValOrZero(sub[3]),
				Ask:             decValOrZero(sub[4]),
				AskSize:         decValOrZero(sub[5]),
				AskPeriod:       int64ValOrZero(sub[6]),
				DailyChange:     decValOrZero(sub[7]),
				DailyChangePerc: decValOrZero(sub[8]),
				LastPrice:       decValOrZero(sub[9]),
				Volume:          decValOrZero(sub[10]),
				High:            decValOrZero(sub[11]),
				Low:             decValOrZero(sub[12]),
			}
		}
		tickers[i] = entry
	}
	return &bitfinex.TickerSnapshot{Snapshot: tickers}, nil
}
