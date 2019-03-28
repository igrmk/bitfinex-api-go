package websocket

import (
	"hash/crc32"
	"sort"
	"strings"
	"sync"

	bitfinex "github.com/igrmk/bitfinex-api-go/v2"
)

type Checksum struct {
	Symbol string
	Value  uint32
}

type Orderbook struct {
	lock sync.Mutex

	symbol string
	bids   []*bitfinex.BookUpdate
	asks   []*bitfinex.BookUpdate
}

func (ob *Orderbook) SetWithSnapshot(bs *bitfinex.BookUpdateSnapshot) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	ob.bids = make([]*bitfinex.BookUpdate, 0)
	ob.asks = make([]*bitfinex.BookUpdate, 0)
	for _, order := range bs.Snapshot {
		if order.Side == bitfinex.Bid {
			ob.bids = append(ob.bids, order)
		} else {
			ob.asks = append(ob.asks, order)
		}
	}
}

func (ob *Orderbook) UpdateWith(bu *bitfinex.BookUpdate) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	side := &ob.asks
	if bu.Side == bitfinex.Bid {
		side = &ob.bids
	}

	// check if first in book
	if len(*side) == 0 {
		*side = append(*side, bu)
		return
	}

	// match price level
	for index, sOrder := range *side {
		if sOrder.Price.Cmp(bu.Price) == 0 {
			if index+1 > len(*(side)) {
				return
			}
			if bu.Count <= 0 {
				// delete if count is equal to zero
				*side = append((*side)[:index], (*side)[index+1:]...)
				return
			} else {
				// remove now and we will add in the code below
				*side = append((*side)[:index], (*side)[index+1:]...)
			}
		}
	}
	*side = append(*side, bu)
	// add to the orderbook and sort lowest to highest
	sort.Slice(*side, func(i, j int) bool {
		if bu.Side == bitfinex.Ask {
			return (*side)[i].Price.Cmp((*side)[j].Price) < 0
		} else {
			return (*side)[i].Price.Cmp((*side)[j].Price) > 0
		}
	})
}

func (ob *Orderbook) Checksum() uint32 {
	ob.lock.Lock()
	defer ob.lock.Unlock()
	var checksumItems []string
	for i := 0; i < 25; i++ {
		if len(ob.bids) > i {
			// append bid
			price := ob.bids[i].Price.String()
			amount := ob.bids[i].Amount.String()
			checksumItems = append(checksumItems, price)
			checksumItems = append(checksumItems, amount)
		}
		if len(ob.asks) > i {
			// append ask
			price := ob.asks[i].Price.String()
			amount := ob.asks[i].Amount.Neg().String()
			checksumItems = append(checksumItems, price)
			checksumItems = append(checksumItems, amount)
		}
	}
	checksumStrings := strings.Join(checksumItems, ":")
	return crc32.ChecksumIEEE([]byte(checksumStrings))
}
