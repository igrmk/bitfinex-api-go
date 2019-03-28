package rest

import (
	"encoding/json"
	"net/url"
	"path"
	"strconv"

	"github.com/igrmk/bitfinex-api-go/v2"
)

type BookService struct {
	Synchronous
}

func (b *BookService) All(symbol string, precision bitfinex.BookPrecision, priceLevels int) (*bitfinex.BookUpdateSnapshot, error) {
	req := NewRequestWithMethod(path.Join("book", symbol, string(precision)), "GET")
	req.Params = make(url.Values)
	req.Params.Add("len", strconv.Itoa(priceLevels))
	raw, err := b.Request(req)

	if err != nil {
		return nil, err
	}

	data := make([][]json.Number, 0, len(raw))
	for _, ifacearr := range raw {
		if arr, ok := ifacearr.([]interface{}); ok {
			sub := make([]json.Number, 0, len(arr))
			for _, iface := range arr {
				if flt, ok := iface.(json.Number); ok {
					sub = append(sub, flt)
				}
			}
			data = append(data, sub)
		}
	}

	book, err := bitfinex.NewBookUpdateSnapshotFromRaw(symbol, string(precision), data)
	if err != nil {
		return nil, err
	}

	return book, nil
}
