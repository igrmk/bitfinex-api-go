package bitfinex

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/igrmk/decimal"
)

// Prefixes for available pairs
const (
	FundingPrefix = "f"
	TradingPrefix = "t"
)

var (
	ErrNotFound = errors.New("not found")
)

// Candle resolutions
const (
	OneMinute      CandleResolution = "1m"
	FiveMinutes    CandleResolution = "5m"
	FifteenMinutes CandleResolution = "15m"
	ThirtyMinutes  CandleResolution = "30m"
	OneHour        CandleResolution = "1h"
	ThreeHours     CandleResolution = "3h"
	SixHours       CandleResolution = "6h"
	TwelveHours    CandleResolution = "12h"
	OneDay         CandleResolution = "1D"
	OneWeek        CandleResolution = "7D"
	TwoWeeks       CandleResolution = "14D"
	OneMonth       CandleResolution = "1M"
)

type Mts int64
type SortOrder int

const (
	OldestFirst SortOrder = 1
	NewestFirst SortOrder = -1
)

type QueryLimit int

const QueryLimitMax QueryLimit = 1000

func CandleResolutionFromString(str string) (CandleResolution, error) {
	switch str {
	case string(OneMinute):
		return OneMinute, nil
	case string(FiveMinutes):
		return FiveMinutes, nil
	case string(FifteenMinutes):
		return FifteenMinutes, nil
	case string(ThirtyMinutes):
		return ThirtyMinutes, nil
	case string(OneHour):
		return OneHour, nil
	case string(ThreeHours):
		return ThreeHours, nil
	case string(SixHours):
		return SixHours, nil
	case string(TwelveHours):
		return TwelveHours, nil
	case string(OneDay):
		return OneDay, nil
	case string(OneWeek):
		return OneWeek, nil
	case string(TwoWeeks):
		return TwoWeeks, nil
	case string(OneMonth):
		return OneMonth, nil
	}
	return OneMinute, fmt.Errorf("could not convert string to resolution: %s", str)
}

// private type--cannot instantiate.
type candleResolution string

// CandleResolution provides a typed set of resolutions for candle subscriptions.
type CandleResolution candleResolution

// Order sides
const (
	Bid OrderSide = 1
	Ask OrderSide = 2
)

// Settings flags

const (
	Dec_s     int = 9
	Time_s    int = 32
	Timestamp int = 32768
	Seq_all   int = 65536
	Checksum  int = 131072
)

type orderSide byte

// OrderSide provides a typed set of order sides.
type OrderSide orderSide

// Book precision levels
const (
	// Aggregate precision levels
	Precision0 BookPrecision = "P0"
	Precision2 BookPrecision = "P2"
	Precision1 BookPrecision = "P1"
	Precision3 BookPrecision = "P3"
	// Raw precision
	PrecisionRawBook BookPrecision = "R0"
)

// private type
type bookPrecision string

// BookPrecision provides a typed book precision level.
type BookPrecision bookPrecision

const (
	// FrequencyRealtime book frequency gives updates as they occur in real-time.
	FrequencyRealtime BookFrequency = "F0"
	// FrequencyTwoPerSecond delivers two book updates per second.
	FrequencyTwoPerSecond BookFrequency = "F1"
	// PriceLevelDefault provides a constant default price level for book subscriptions.
	PriceLevelDefault int = 25
)

type bookFrequency string

// BookFrequency provides a typed book frequency.
type BookFrequency bookFrequency

const (
	OrderFlagHidden   int = 64
	OrderFlagClose    int = 512
	OrderFlagPostOnly int = 4096
	OrderFlagOCO      int = 16384
)

// OrderNewRequest represents an order to be posted to the bitfinex websocket
// service.
type OrderNewRequest struct {
	GID           int64           `json:"gid"`
	CID           int64           `json:"cid"`
	Type          string          `json:"type"`
	Symbol        string          `json:"symbol"`
	Amount        decimal.Decimal `json:"amount,string"`
	Price         decimal.Decimal `json:"price,string"`
	PriceTrailing decimal.Decimal `json:"price_trailing,string,omitempty"`
	PriceAuxLimit decimal.Decimal `json:"price_aux_limit,string,omitempty"`
	PriceOcoStop  decimal.Decimal `json:"price_oco_stop,string,omitempty"`
	Hidden        bool            `json:"hidden,omitempty"`
	PostOnly      bool            `json:"postonly,omitempty"`
	Close         bool            `json:"close,omitempty"`
	OcoOrder      bool            `json:"oco_order,omitempty"`
	TimeInForce   string          `json:"tif,omitempty"`
}

// MarshalJSON converts the order object into the format required by the bitfinex
// websocket service.
func (o *OrderNewRequest) MarshalJSON() ([]byte, error) {
	aux := struct {
		GID           int64           `json:"gid"`
		CID           int64           `json:"cid"`
		Type          string          `json:"type"`
		Symbol        string          `json:"symbol"`
		Amount        decimal.Decimal `json:"amount,string"`
		Price         decimal.Decimal `json:"price,string"`
		PriceTrailing decimal.Decimal `json:"price_trailing,string,omitempty"`
		PriceAuxLimit decimal.Decimal `json:"price_aux_limit,string,omitempty"`
		PriceOcoStop  decimal.Decimal `json:"price_oco_stop,string,omitempty"`
		TimeInForce   string          `json:"tif,omitempty"`
		Flags         int             `json:"flags,omitempty"`
	}{
		GID:           o.GID,
		CID:           o.CID,
		Type:          o.Type,
		Symbol:        o.Symbol,
		Amount:        o.Amount,
		Price:         o.Price,
		PriceTrailing: o.PriceTrailing,
		PriceAuxLimit: o.PriceAuxLimit,
		PriceOcoStop:  o.PriceOcoStop,
	}

	if o.Hidden {
		aux.Flags = aux.Flags + OrderFlagHidden
	}

	if o.PostOnly {
		aux.Flags = aux.Flags + OrderFlagPostOnly
	}

	if o.OcoOrder {
		aux.Flags = aux.Flags + OrderFlagOCO
	}

	if o.Close {
		aux.Flags = aux.Flags + OrderFlagClose
	}

	body := []interface{}{0, "on", nil, aux}
	return json.Marshal(&body)
}

type OrderUpdateRequest struct {
	ID            int64   `json:"id"`
	GID           int64   `json:"gid,omitempty"`
	Price         float64 `json:"price,string,omitempty"`
	Amount        float64 `json:"amount,string,omitempty"`
	Delta         float64 `json:"delta,string,omitempty"`
	PriceTrailing float64 `json:"price_trailing,string,omitempty"`
	PriceAuxLimit float64 `json:"price_aux_limit,string,omitempty"`
	Hidden        bool    `json:"hidden,omitempty"`
	PostOnly      bool    `json:"postonly,omitempty"`
}

// MarshalJSON converts the order object into the format required by the bitfinex
// websocket service.
func (o *OrderUpdateRequest) MarshalJSON() ([]byte, error) {
	aux := struct {
		ID            int64   `json:"id"`
		GID           int64   `json:"gid,omitempty"`
		Price         float64 `json:"price,string,omitempty"`
		Amount        float64 `json:"amount,string,omitempty"`
		Delta         float64 `json:"delta,string,omitempty"`
		PriceTrailing float64 `json:"price_trailing,string,omitempty"`
		PriceAuxLimit float64 `json:"price_aux_limit,string,omitempty"`
		Hidden        bool    `json:"hidden,omitempty"`
		PostOnly      bool    `json:"postonly,omitempty"`
		Flags         int     `json:"flags,omitempty"`
	}{
		ID:            o.ID,
		GID:           o.GID,
		Amount:        o.Amount,
		Price:         o.Price,
		PriceTrailing: o.PriceTrailing,
		PriceAuxLimit: o.PriceAuxLimit,
		Delta:         o.Delta,
	}

	if o.Hidden {
		aux.Flags = aux.Flags + OrderFlagHidden
	}

	if o.PostOnly {
		aux.Flags = aux.Flags + OrderFlagPostOnly
	}

	body := []interface{}{0, "ou", nil, aux}
	return json.Marshal(&body)
}

// OrderCancelRequest represents an order cancel request.
// An order can be cancelled using the internal ID or a
// combination of Client ID (CID) and the daten for the given
// CID.
type OrderCancelRequest struct {
	ID      int64  `json:"id,omitempty"`
	CID     int64  `json:"cid,omitempty"`
	CIDDate string `json:"cid_date,omitempty"`
}

// MarshalJSON converts the order cancel object into the format required by the
// bitfinex websocket service.
func (o *OrderCancelRequest) MarshalJSON() ([]byte, error) {
	aux := struct {
		ID      int64  `json:"id,omitempty"`
		CID     int64  `json:"cid,omitempty"`
		CIDDate string `json:"cid_date,omitempty"`
	}{
		ID:      o.ID,
		CID:     o.CID,
		CIDDate: o.CIDDate,
	}

	body := []interface{}{0, "oc", nil, aux}
	return json.Marshal(&body)
}

// TODO: MultiOrderCancelRequest represents an order cancel request.

type Heartbeat struct {
	//ChannelIDs []int64
}

// OrderType represents the types orders the bitfinex platform can handle.
type OrderType string

const (
	OrderTypeMarket               = "MARKET"
	OrderTypeExchangeMarket       = "EXCHANGE MARKET"
	OrderTypeLimit                = "LIMIT"
	OrderTypeExchangeLimit        = "EXCHANGE LIMIT"
	OrderTypeStop                 = "STOP"
	OrderTypeExchangeStop         = "EXCHANGE STOP"
	OrderTypeTrailingStop         = "TRAILING STOP"
	OrderTypeExchangeTrailingStop = "EXCHANGE TRAILING STOP"
	OrderTypeFOK                  = "FOK"
	OrderTypeExchangeFOK          = "EXCHANGE FOK"
	OrderTypeStopLimit            = "STOP LIMIT"
	OrderTypeExchangeStopLimit    = "EXCHANGE STOP LIMIT"
)

// OrderStatus represents the possible statuses an order can be in.
type OrderStatus string

const (
	OrderStatusActive          OrderStatus = "ACTIVE"
	OrderStatusExecuted        OrderStatus = "EXECUTED"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
)

// Order as returned from the bitfinex websocket service.
type Order struct {
	ID            int64
	GID           int64
	CID           int64
	Symbol        string
	MTSCreated    int64
	MTSUpdated    int64
	Amount        decimal.Decimal
	AmountOrig    decimal.Decimal
	Type          string
	TypePrev      string
	Flags         int64
	Status        OrderStatus
	Price         decimal.Decimal
	PriceAvg      decimal.Decimal
	PriceTrailing decimal.Decimal
	PriceAuxLimit decimal.Decimal
	Notify        bool
	Hidden        bool
	PlacedID      int64
}

// NewOrderFromRaw takes the raw list of values as returned from the websocket
// service and tries to convert it into an Order.
func NewOrderFromRaw(raw []interface{}) (o *Order, err error) {
	if len(raw) == 12 {
		o = &Order{
			ID:         Int64ValOrZero(raw[0]),
			Symbol:     strValOrEmpty(raw[1]),
			Amount:     DecValOrZero(raw[2]),
			AmountOrig: DecValOrZero(raw[3]),
			Type:       strValOrEmpty(raw[4]),
			Status:     OrderStatus(strValOrEmpty(raw[5])),
			Price:      DecValOrZero(raw[6]),
			PriceAvg:   DecValOrZero(raw[7]),
			MTSUpdated: Int64ValOrZero(raw[8]),
			// 3 trailing zeroes, what do they map to?
		}
	} else if len(raw) < 26 {
		return o, fmt.Errorf("data slice too short for order: %#v", raw)
	} else {
		// TODO: API docs say ID, GID, CID, MTS_CREATE, MTS_UPDATE are int but API returns float
		o = &Order{
			ID:            Int64ValOrZero(raw[0]),
			GID:           Int64ValOrZero(raw[1]),
			CID:           Int64ValOrZero(raw[2]),
			Symbol:        strValOrEmpty(raw[3]),
			MTSCreated:    Int64ValOrZero(raw[4]),
			MTSUpdated:    Int64ValOrZero(raw[5]),
			Amount:        DecValOrZero(raw[6]),
			AmountOrig:    DecValOrZero(raw[7]),
			Type:          strValOrEmpty(raw[8]),
			TypePrev:      strValOrEmpty(raw[9]),
			Flags:         Int64ValOrZero(raw[12]),
			Status:        OrderStatus(strValOrEmpty(raw[13])),
			Price:         DecValOrZero(raw[16]),
			PriceAvg:      DecValOrZero(raw[17]),
			PriceTrailing: DecValOrZero(raw[18]),
			PriceAuxLimit: DecValOrZero(raw[19]),
			Notify:        bValOrFalse(raw[23]),
			Hidden:        bValOrFalse(raw[24]),
			PlacedID:      Int64ValOrZero(raw[25]),
		}
	}

	return
}

// OrderSnapshotFromRaw takes a raw list of values as returned from the websocket
// service and tries to convert it into an OrderSnapshot.
func NewOrderSnapshotFromRaw(raw []interface{}) (s *OrderSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	os := make([]*Order, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := NewOrderFromRaw(l)
				if err != nil {
					return s, err
				}
				os = append(os, o)
			}
		}
	default:
		return s, fmt.Errorf("not an order snapshot")
	}
	s = &OrderSnapshot{Snapshot: os}

	return
}

// OrderSnapshot is a collection of Orders that would usually be sent on
// inital connection.
type OrderSnapshot struct {
	Snapshot []*Order
}

// OrderUpdate is an Order that gets sent out after every change to an
// order.
type OrderUpdate Order

// OrderNew gets sent out after an Order was created successfully.
type OrderNew Order

// OrderCancel gets sent out after an Order was cancelled successfully.
type OrderCancel Order

type PositionStatus string

const (
	PositionStatusActive PositionStatus = "ACTIVE"
	PositionStatusClosed PositionStatus = "CLOSED"
)

type Position struct {
	Symbol               string
	Status               PositionStatus
	Amount               decimal.Decimal
	BasePrice            decimal.Decimal
	MarginFunding        decimal.Decimal
	MarginFundingType    int64
	ProfitLoss           decimal.Decimal
	ProfitLossPercentage decimal.Decimal
	LiquidationPrice     decimal.Decimal
	Leverage             decimal.Decimal
}

func NewPositionFromRaw(raw []interface{}) (o *Position, err error) {
	if len(raw) == 6 {
		o = &Position{
			Symbol:            strValOrEmpty(raw[0]),
			Status:            PositionStatus(strValOrEmpty(raw[1])),
			Amount:            DecValOrZero(raw[2]),
			BasePrice:         DecValOrZero(raw[3]),
			MarginFunding:     DecValOrZero(raw[4]),
			MarginFundingType: Int64ValOrZero(raw[5]),
		}
	} else if len(raw) < 10 {
		return o, fmt.Errorf("data slice too short for position: %#v", raw)
	} else {
		o = &Position{
			Symbol:               strValOrEmpty(raw[0]),
			Status:               PositionStatus(strValOrEmpty(raw[1])),
			Amount:               DecValOrZero(raw[2]),
			BasePrice:            DecValOrZero(raw[3]),
			MarginFunding:        DecValOrZero(raw[4]),
			MarginFundingType:    Int64ValOrZero(raw[5]),
			ProfitLoss:           DecValOrZero(raw[6]),
			ProfitLossPercentage: DecValOrZero(raw[7]),
			LiquidationPrice:     DecValOrZero(raw[8]),
			Leverage:             DecValOrZero(raw[9]),
		}
	}
	return
}

type PositionSnapshot struct {
	Snapshot []*Position
}
type PositionNew Position
type PositionUpdate Position
type PositionCancel Position

func NewPositionSnapshotFromRaw(raw []interface{}) (s *PositionSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	ps := make([]*Position, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				p, err := NewPositionFromRaw(l)
				if err != nil {
					return s, err
				}
				ps = append(ps, p)
			}
		}
	default:
		return s, fmt.Errorf("not a position snapshot")
	}
	s = &PositionSnapshot{Snapshot: ps}

	return
}

// Trade represents a trade on the public data feed.
type Trade struct {
	Pair   string
	ID     int64
	MTS    int64
	Amount decimal.Decimal
	Price  decimal.Decimal
	Side   OrderSide
}

func NewTradeFromRaw(pair string, raw []interface{}) (o *Trade, err error) {
	if len(raw) < 4 {
		return o, fmt.Errorf("data slice too short for trade: %#v", raw)
	}

	amt := DecValOrZero(raw[2])
	var side OrderSide
	if amt.Cmp(decimal.Zero) > 0 {
		side = Bid
	} else {
		side = Ask
	}

	o = &Trade{
		Pair:   pair,
		ID:     Int64ValOrZero(raw[0]),
		MTS:    Int64ValOrZero(raw[1]),
		Amount: amt.Abs(),
		Price:  DecValOrZero(raw[3]),
		Side:   side,
	}

	return
}

type TradeSnapshot struct {
	Snapshot []*Trade
}

func NewTradeSnapshotFromRaw(pair string, raw [][]json.Number) (*TradeSnapshot, error) {
	if len(raw) <= 0 {
		return nil, fmt.Errorf("data slice is too short for trade snapshot: %#v", raw)
	}
	snapshot := make([]*Trade, 0)
	for _, flt := range raw {
		t, err := NewTradeFromRaw(pair, toInterfaceSlice(flt))
		if err == nil {
			snapshot = append(snapshot, t)
		}
	}

	return &TradeSnapshot{Snapshot: snapshot}, nil
}

// TradeExecutionUpdate represents a full update to a trade on the private data feed.  Following a TradeExecution,
// TradeExecutionUpdates include additional details, e.g. the trade's execution ID (TradeID).
type TradeExecutionUpdate struct {
	ID          int64
	Pair        string
	MTS         int64
	OrderID     int64
	ExecAmount  decimal.Decimal
	ExecPrice   decimal.Decimal
	OrderType   string
	OrderPrice  decimal.Decimal
	Maker       int
	Fee         decimal.Decimal
	FeeCurrency string
}

// public trade update just looks like a trade
func NewTradeExecutionUpdateFromRaw(raw []interface{}) (o *TradeExecutionUpdate, err error) {
	if len(raw) == 4 {
		o = &TradeExecutionUpdate{
			ID:         Int64ValOrZero(raw[0]),
			MTS:        Int64ValOrZero(raw[1]),
			ExecAmount: DecValOrZero(raw[2]),
			ExecPrice:  DecValOrZero(raw[3]),
		}
		return
	}
	if len(raw) == 11 {
		o = &TradeExecutionUpdate{
			ID:          Int64ValOrZero(raw[0]),
			Pair:        strValOrEmpty(raw[1]),
			MTS:         Int64ValOrZero(raw[2]),
			OrderID:     Int64ValOrZero(raw[3]),
			ExecAmount:  DecValOrZero(raw[4]),
			ExecPrice:   DecValOrZero(raw[5]),
			OrderType:   strValOrEmpty(raw[6]),
			OrderPrice:  DecValOrZero(raw[7]),
			Maker:       int(Int64ValOrZero(raw[8])),
			Fee:         DecValOrZero(raw[9]),
			FeeCurrency: strValOrEmpty(raw[10]),
		}
		return
	}
	return o, fmt.Errorf("data slice too short for trade update: %#v", raw)
}

type TradeExecutionUpdateSnapshot struct {
	Snapshot []*TradeExecutionUpdate
}
type HistoricalTradeSnapshot TradeExecutionUpdateSnapshot

func NewTradeExecutionUpdateSnapshotFromRaw(raw []interface{}) (s *TradeExecutionUpdateSnapshot, err error) {
	if len(raw) == 0 {
		return
	}
	ts := make([]*TradeExecutionUpdate, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				t, err := NewTradeExecutionUpdateFromRaw(l)
				if err != nil {
					return s, err
				}
				ts = append(ts, t)
			}
		}
	default:
		return s, fmt.Errorf("not a trade snapshot: %#v", raw)
	}
	s = &TradeExecutionUpdateSnapshot{Snapshot: ts}

	return
}

// TradeExecution represents the first message receievd for a trade on the private data feed.
type TradeExecution struct {
	ID         int64
	Pair       string
	MTS        int64
	OrderID    int64
	Amount     decimal.Decimal
	Price      decimal.Decimal
	OrderType  string
	OrderPrice decimal.Decimal
	Maker      int
}

func NewTradeExecutionFromRaw(raw []interface{}) (o *TradeExecution, err error) {
	if len(raw) < 6 {
		log.Printf("[ERROR] not enough members (%d, need at least 6) for trade execution: %#v", len(raw), raw)
		return o, fmt.Errorf("data slice too short for trade execution: %#v", raw)
	}

	// trade executions sometimes omit order type, price, and maker flag
	o = &TradeExecution{
		ID:      Int64ValOrZero(raw[0]),
		Pair:    strValOrEmpty(raw[1]),
		MTS:     Int64ValOrZero(raw[2]),
		OrderID: Int64ValOrZero(raw[3]),
		Amount:  DecValOrZero(raw[4]),
		Price:   DecValOrZero(raw[5]),
	}

	if len(raw) >= 9 {
		o.OrderType = strValOrEmpty(raw[6])
		o.OrderPrice = DecValOrZero(raw[7])
		o.Maker = int(Int64ValOrZero(raw[8]))
	}

	return
}

type Wallet struct {
	Type              string
	Currency          string
	Balance           decimal.Decimal
	UnsettledInterest decimal.Decimal
	BalanceAvailable  decimal.Decimal
}

func NewWalletFromRaw(raw []interface{}) (o *Wallet, err error) {
	if len(raw) == 4 {
		o = &Wallet{
			Type:              strValOrEmpty(raw[0]),
			Currency:          strValOrEmpty(raw[1]),
			Balance:           DecValOrZero(raw[2]),
			UnsettledInterest: DecValOrZero(raw[3]),
		}
	} else if len(raw) < 5 {
		return o, fmt.Errorf("data slice too short for wallet: %#v", raw)
	} else {
		o = &Wallet{
			Type:              strValOrEmpty(raw[0]),
			Currency:          strValOrEmpty(raw[1]),
			Balance:           DecValOrZero(raw[2]),
			UnsettledInterest: DecValOrZero(raw[3]),
			BalanceAvailable:  DecValOrZero(raw[4]),
		}
	}
	return
}

type WalletUpdate Wallet
type WalletSnapshot struct {
	Snapshot []*Wallet
}

func NewWalletSnapshotFromRaw(raw []interface{}) (s *WalletSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	ws := make([]*Wallet, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := NewWalletFromRaw(l)
				if err != nil {
					return s, err
				}
				ws = append(ws, o)
			}
		}
	default:
		return s, fmt.Errorf("not an wallet snapshot")
	}
	s = &WalletSnapshot{Snapshot: ws}

	return
}

type BalanceInfo struct {
	TotalAUM decimal.Decimal
	NetAUM   decimal.Decimal
	/*WalletType string
	Currency   string*/
}

func NewBalanceInfoFromRaw(raw []interface{}) (o *BalanceInfo, err error) {
	if len(raw) < 2 {
		return o, fmt.Errorf("data slice too short for balance info: %#v", raw)
	}

	o = &BalanceInfo{
		TotalAUM: DecValOrZero(raw[0]),
		NetAUM:   DecValOrZero(raw[1]),
		/*WalletType: sValOrEmpty(raw[2]),
		Currency:   sValOrEmpty(raw[3]),*/
	}

	return
}

type BalanceUpdate BalanceInfo

// marginInfoFromRaw returns either a MarginInfoBase or MarginInfoUpdate, since
// the Margin Info is split up into a base and per symbol parts.
func NewMarginInfoFromRaw(raw []interface{}) (o interface{}, err error) {
	if len(raw) < 2 {
		return o, fmt.Errorf("data slice too short for margin info base: %#v", raw)
	}

	typ, ok := raw[0].(string)
	if !ok {
		return o, fmt.Errorf("expected margin info type in first position for margin info but got %#v", raw)
	}

	if len(raw) == 2 && typ == "base" { // This should be ["base", [...]]
		data, ok := raw[1].([]interface{})
		if !ok {
			return o, fmt.Errorf("expected margin info array in second position for margin info but got %#v", raw)
		}

		return NewMarginInfoBaseFromRaw(data)
	} else if len(raw) == 3 && typ == "sym" { // This should be ["sym", SYMBOL, [...]]
		symbol, ok := raw[1].(string)
		if !ok {
			return o, fmt.Errorf("expected margin info symbol in second position for margin info update but got %#v", raw)
		}

		data, ok := raw[2].([]interface{})
		if !ok {
			return o, fmt.Errorf("expected margin info array in third position for margin info update but got %#v", raw)
		}

		return NewMarginInfoUpdateFromRaw(symbol, data)
	}

	return nil, fmt.Errorf("invalid margin info type in %#v", raw)
}

type MarginInfoUpdate struct {
	Symbol          string
	TradableBalance decimal.Decimal
}

func NewMarginInfoUpdateFromRaw(symbol string, raw []interface{}) (o *MarginInfoUpdate, err error) {
	if len(raw) < 1 {
		return o, fmt.Errorf("data slice too short for margin info update: %#v", raw)
	}

	o = &MarginInfoUpdate{
		Symbol:          symbol,
		TradableBalance: DecValOrZero(raw[0]),
	}

	return
}

type MarginInfoBase struct {
	UserProfitLoss decimal.Decimal
	UserSwaps      decimal.Decimal
	MarginBalance  decimal.Decimal
	MarginNet      decimal.Decimal
}

func NewMarginInfoBaseFromRaw(raw []interface{}) (o *MarginInfoBase, err error) {
	if len(raw) < 4 {
		return o, fmt.Errorf("data slice too short for margin info base: %#v", raw)
	}

	o = &MarginInfoBase{
		UserProfitLoss: DecValOrZero(raw[0]),
		UserSwaps:      DecValOrZero(raw[1]),
		MarginBalance:  DecValOrZero(raw[2]),
		MarginNet:      DecValOrZero(raw[3]),
	}

	return
}

type FundingInfo struct {
	Symbol       string
	YieldLoan    decimal.Decimal
	YieldLend    decimal.Decimal
	DurationLoan decimal.Decimal
	DurationLend decimal.Decimal
}

func NewFundingInfoFromRaw(raw []interface{}) (o *FundingInfo, err error) {
	if len(raw) < 3 { // "sym", symbol, data
		return o, fmt.Errorf("data slice too short for funding info: %#v", raw)
	}

	sym, ok := raw[1].(string)
	if !ok {
		return o, fmt.Errorf("expected symbol in second position of funding info: %v", raw)
	}

	data, ok := raw[2].([]interface{})
	if !ok {
		return o, fmt.Errorf("expected list in third position of funding info: %v", raw)
	}

	if len(data) < 4 {
		return o, fmt.Errorf("data too short: %#v", data)
	}

	o = &FundingInfo{
		Symbol:       sym,
		YieldLoan:    DecValOrZero(data[0]),
		YieldLend:    DecValOrZero(data[1]),
		DurationLoan: DecValOrZero(data[2]),
		DurationLend: DecValOrZero(data[3]),
	}

	return
}

type OfferStatus string

const (
	OfferStatusActive          OfferStatus = "ACTIVE"
	OfferStatusExecuted        OfferStatus = "EXECUTED"
	OfferStatusPartiallyFilled OfferStatus = "PARTIALLY FILLED"
	OfferStatusCanceled        OfferStatus = "CANCELED"
)

type Offer struct {
	ID         int64
	Symbol     string
	MTSCreated int64
	MTSUpdated int64
	Amount     decimal.Decimal
	AmountOrig decimal.Decimal
	Type       string
	Flags      interface{}
	Status     OfferStatus
	Rate       decimal.Decimal
	Period     int64
	Notify     bool
	Hidden     bool
	Insure     bool
	Renew      bool
	RateReal   decimal.Decimal
}

func NewOfferFromRaw(raw []interface{}) (o *Offer, err error) {
	if len(raw) < 21 {
		return o, fmt.Errorf("data slice too short for offer: %#v", raw)
	}

	o = &Offer{
		ID:         Int64ValOrZero(raw[0]),
		Symbol:     strValOrEmpty(raw[1]),
		MTSCreated: Int64ValOrZero(raw[2]),
		MTSUpdated: Int64ValOrZero(raw[3]),
		Amount:     DecValOrZero(raw[4]),
		AmountOrig: DecValOrZero(raw[5]),
		Type:       strValOrEmpty(raw[6]),
		Flags:      raw[9],
		Status:     OfferStatus(strValOrEmpty(raw[10])),
		Rate:       DecValOrZero(raw[14]),
		Period:     Int64ValOrZero(raw[15]),
		Notify:     bValOrFalse(raw[16]),
		Hidden:     bValOrFalse(raw[17]),
		Insure:     bValOrFalse(raw[18]),
		Renew:      bValOrFalse(raw[19]),
		RateReal:   DecValOrZero(raw[20]),
	}

	return
}

type FundingOfferNew Offer
type FundingOfferUpdate Offer
type FundingOfferCancel Offer
type FundingOfferSnapshot struct {
	Snapshot []*Offer
}

func NewFundingOfferSnapshotFromRaw(raw []interface{}) (snap *FundingOfferSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	fos := make([]*Offer, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := NewOfferFromRaw(l)
				if err != nil {
					return snap, err
				}
				fos = append(fos, o)
			}
		}
	default:
		return snap, fmt.Errorf("not a funding offer snapshot")
	}

	snap = &FundingOfferSnapshot{
		Snapshot: fos,
	}

	return
}

type HistoricalOffer Offer

type CreditStatus string

const (
	CreditStatusActive          CreditStatus = "ACTIVE"
	CreditStatusExecuted        CreditStatus = "EXECUTED"
	CreditStatusPartiallyFilled CreditStatus = "PARTIALLY FILLED"
	CreditStatusCanceled        CreditStatus = "CANCELED"
)

type Credit struct {
	ID            int64
	Symbol        string
	Side          string
	MTSCreated    int64
	MTSUpdated    int64
	Amout         decimal.Decimal
	Flags         interface{}
	Status        CreditStatus
	Rate          decimal.Decimal
	Period        int64
	MTSOpened     int64
	MTSLastPayout int64
	Notify        bool
	Hidden        bool
	Insure        bool
	Renew         bool
	RateReal      decimal.Decimal
	NoClose       bool
	PositionPair  string
}

func NewCreditFromRaw(raw []interface{}) (o *Credit, err error) {
	if len(raw) < 22 {
		return o, fmt.Errorf("data slice too short for offer: %#v", raw)
	}

	o = &Credit{
		ID:            Int64ValOrZero(raw[0]),
		Symbol:        strValOrEmpty(raw[1]),
		Side:          strValOrEmpty(raw[2]),
		MTSCreated:    Int64ValOrZero(raw[3]),
		MTSUpdated:    Int64ValOrZero(raw[4]),
		Amout:         DecValOrZero(raw[5]),
		Flags:         raw[6],
		Status:        CreditStatus(strValOrEmpty(raw[7])),
		Rate:          DecValOrZero(raw[11]),
		Period:        Int64ValOrZero(raw[12]),
		MTSOpened:     Int64ValOrZero(raw[13]),
		MTSLastPayout: Int64ValOrZero(raw[14]),
		Notify:        bValOrFalse(raw[15]),
		Hidden:        bValOrFalse(raw[16]),
		Insure:        bValOrFalse(raw[17]),
		Renew:         bValOrFalse(raw[18]),
		RateReal:      DecValOrZero(raw[19]),
		NoClose:       bValOrFalse(raw[20]),
		PositionPair:  strValOrEmpty(raw[21]),
	}

	return
}

type HistoricalCredit Credit
type FundingCreditNew Credit
type FundingCreditUpdate Credit
type FundingCreditCancel Credit

type FundingCreditSnapshot struct {
	Snapshot []*Credit
}

func NewFundingCreditSnapshotFromRaw(raw []interface{}) (snap *FundingCreditSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	fcs := make([]*Credit, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := NewCreditFromRaw(l)
				if err != nil {
					return snap, err
				}
				fcs = append(fcs, o)
			}
		}
	default:
		return snap, fmt.Errorf("not a funding credit snapshot")
	}
	snap = &FundingCreditSnapshot{
		Snapshot: fcs,
	}

	return
}

type LoanStatus string

const (
	LoanStatusActive          LoanStatus = "ACTIVE"
	LoanStatusExecuted        LoanStatus = "EXECUTED"
	LoanStatusPartiallyFilled LoanStatus = "PARTIALLY FILLED"
	LoanStatusCanceled        LoanStatus = "CANCELED"
)

type Loan struct {
	ID            int64
	Symbol        string
	Side          string
	MTSCreated    int64
	MTSUpdated    int64
	Amout         decimal.Decimal
	Flags         interface{}
	Status        LoanStatus
	Rate          decimal.Decimal
	Period        int64
	MTSOpened     int64
	MTSLastPayout int64
	Notify        bool
	Hidden        bool
	Insure        bool
	Renew         bool
	RateReal      decimal.Decimal
	NoClose       bool
}

func NewLoanFromRaw(raw []interface{}) (o *Loan, err error) {
	if len(raw) < 21 {
		return o, fmt.Errorf("data slice too short (len=%d) for loan: %#v", len(raw), raw)
	}

	o = &Loan{
		ID:            Int64ValOrZero(raw[0]),
		Symbol:        strValOrEmpty(raw[1]),
		Side:          strValOrEmpty(raw[2]),
		MTSCreated:    Int64ValOrZero(raw[3]),
		MTSUpdated:    Int64ValOrZero(raw[4]),
		Amout:         DecValOrZero(raw[5]),
		Flags:         raw[6],
		Status:        LoanStatus(strValOrEmpty(raw[7])),
		Rate:          DecValOrZero(raw[11]),
		Period:        Int64ValOrZero(raw[12]),
		MTSOpened:     Int64ValOrZero(raw[13]),
		MTSLastPayout: Int64ValOrZero(raw[14]),
		Notify:        bValOrFalse(raw[15]),
		Hidden:        bValOrFalse(raw[16]),
		Insure:        bValOrFalse(raw[17]),
		Renew:         bValOrFalse(raw[18]),
		RateReal:      DecValOrZero(raw[19]),
		NoClose:       bValOrFalse(raw[20]),
	}

	return o, nil
}

type HistoricalLoan Loan
type FundingLoanNew Loan
type FundingLoanUpdate Loan
type FundingLoanCancel Loan

type FundingLoanSnapshot struct {
	Snapshot []*Loan
}

func NewFundingLoanSnapshotFromRaw(raw []interface{}) (snap *FundingLoanSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	fls := make([]*Loan, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := NewLoanFromRaw(l)
				if err != nil {
					return snap, err
				}
				fls = append(fls, o)
			}
		}
	default:
		return snap, fmt.Errorf("not a funding loan snapshot")
	}
	snap = &FundingLoanSnapshot{
		Snapshot: fls,
	}

	return
}

type FundingTrade struct {
	ID         int64
	Symbol     string
	MTSCreated int64
	OfferID    int64
	Amount     decimal.Decimal
	Rate       decimal.Decimal
	Period     int64
	Maker      int64
}

func NewFundingTradeFromRaw(raw []interface{}) (o *FundingTrade, err error) {
	if len(raw) < 8 {
		return o, fmt.Errorf("data slice too short for funding trade: %#v", raw)
	}

	o = &FundingTrade{
		ID:         Int64ValOrZero(raw[0]),
		Symbol:     strValOrEmpty(raw[1]),
		MTSCreated: Int64ValOrZero(raw[2]),
		OfferID:    Int64ValOrZero(raw[3]),
		Amount:     DecValOrZero(raw[4]),
		Rate:       DecValOrZero(raw[5]),
		Period:     Int64ValOrZero(raw[6]),
		Maker:      Int64ValOrZero(raw[7]),
	}

	return
}

type FundingTradeExecution FundingTrade
type FundingTradeUpdate FundingTrade
type FundingTradeSnapshot struct {
	Snapshot []*FundingTrade
}
type HistoricalFundingTradeSnapshot FundingTradeSnapshot

func NewFundingTradeSnapshotFromRaw(raw []interface{}) (snap *FundingTradeSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	fts := make([]*FundingTrade, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := NewFundingTradeFromRaw(l)
				if err != nil {
					return snap, err
				}
				fts = append(fts, o)
			}
		}
	default:
		return snap, fmt.Errorf("not a funding trade snapshot")
	}
	snap = &FundingTradeSnapshot{
		Snapshot: fts,
	}

	return
}

type Notification struct {
	MTS        int64
	Type       string
	MessageID  int64
	NotifyInfo interface{}
	Code       int64
	Status     string
	Text       string
}

func NewNotificationFromRaw(raw []interface{}) (o *Notification, err error) {
	if len(raw) < 8 {
		return o, fmt.Errorf("data slice too short for notification: %#v", raw)
	}

	o = &Notification{
		MTS:       Int64ValOrZero(raw[0]),
		Type:      strValOrEmpty(raw[1]),
		MessageID: Int64ValOrZero(raw[2]),
		//NotifyInfo: raw[4],
		Code:   Int64ValOrZero(raw[5]),
		Status: strValOrEmpty(raw[6]),
		Text:   strValOrEmpty(raw[7]),
	}

	// raw[4] = notify info
	var nraw []interface{}
	if raw[4] != nil {
		nraw = raw[4].([]interface{})
		switch o.Type {
		case "on-req":
			on, err := NewOrderFromRaw(nraw)
			if err != nil {
				return o, err
			}
			orderNew := OrderNew(*on)
			o.NotifyInfo = &orderNew
		case "oc-req":
			oc, err := NewOrderFromRaw(nraw)
			if err != nil {
				return o, err
			}
			orderCancel := OrderCancel(*oc)
			o.NotifyInfo = &orderCancel
		case "fon-req":
			fon, err := NewOfferFromRaw(nraw)
			if err != nil {
				return o, err
			}
			fundingOffer := FundingOfferNew(*fon)
			o.NotifyInfo = &fundingOffer
		case "foc-req":
			foc, err := NewOfferFromRaw(nraw)
			if err != nil {
				return o, err
			}
			fundingOffer := FundingOfferCancel(*foc)
			o.NotifyInfo = &fundingOffer
		case "uca":
			o.NotifyInfo = raw[4]
		}
	}

	return
}

type Ticker struct {
	Symbol          string
	FRR             decimal.Decimal
	Bid             decimal.Decimal
	BidPeriod       int64
	BidSize         decimal.Decimal
	Ask             decimal.Decimal
	AskPeriod       int64
	AskSize         decimal.Decimal
	DailyChange     decimal.Decimal
	DailyChangePerc decimal.Decimal
	LastPrice       decimal.Decimal
	Volume          decimal.Decimal
	High            decimal.Decimal
	Low             decimal.Decimal
}

type TickerUpdate Ticker
type TickerSnapshot struct {
	Snapshot []*Ticker
}

func NewTickerSnapshotFromRaw(symbol string, raw [][]json.Number) (*TickerSnapshot, error) {
	if len(raw) <= 0 {
		return nil, fmt.Errorf("data slice too short for ticker snapshot: %#v", raw)
	}
	snap := make([]*Ticker, 0)
	for _, f := range raw {
		c, err := NewTickerFromRaw(symbol, toInterfaceSlice(f))
		if err == nil {
			snap = append(snap, c)
		}
	}
	return &TickerSnapshot{Snapshot: snap}, nil
}

func NewTickerFromRaw(symbol string, raw []interface{}) (t *Ticker, err error) {
	if len(raw) < 10 {
		return t, fmt.Errorf("data slice too short for ticker, expected %d got %d: %#v", 10, len(raw), raw)
	}

	t = &Ticker{
		Symbol:          symbol,
		Bid:             DecValOrZero(raw[0]),
		BidSize:         DecValOrZero(raw[1]),
		Ask:             DecValOrZero(raw[2]),
		AskSize:         DecValOrZero(raw[3]),
		DailyChange:     DecValOrZero(raw[4]),
		DailyChangePerc: DecValOrZero(raw[5]),
		LastPrice:       DecValOrZero(raw[6]),
		Volume:          DecValOrZero(raw[7]),
		High:            DecValOrZero(raw[8]),
		Low:             DecValOrZero(raw[9]),
	}

	return
}

type bookAction byte

// BookAction represents a new/update or removal for a book entry.
type BookAction bookAction

const (
	//BookUpdateEntry represents a new or updated book entry.
	BookUpdateEntry BookAction = 0
	//BookRemoveEntry represents a removal of a book entry.
	BookRemoveEntry BookAction = 1
)

// BookUpdate represents an order book price update.
type BookUpdate struct {
	ID     int64           // the book update ID, optional
	Symbol string          // book symbol
	Price  decimal.Decimal // updated price
	Count  int64           // updated count, optional
	Amount decimal.Decimal // updated amount
	Side   OrderSide       // side
	Action BookAction      // action (add/remove)
}

type BookUpdateSnapshot struct {
	Symbol   string        // book symbol
	Snapshot []*BookUpdate // book levels
}

func NewBookUpdateSnapshotFromRaw(symbol, precision string, raw [][]json.Number) (*BookUpdateSnapshot, error) {
	snap := make([]*BookUpdate, len(raw))
	for i, f := range raw {
		b, err := NewBookUpdateFromRaw(symbol, precision, toInterfaceSlice(f))
		if err != nil {
			return nil, err
		}
		snap[i] = b
	}
	return &BookUpdateSnapshot{Symbol: symbol, Snapshot: snap}, nil
}

func IsRawBook(precision string) bool {
	return precision == "R0"
}

// NewBookUpdateFromRaw creates a new book update object from raw data.  Precision determines how
// to interpret the side (baked into Count versus Amount)
// raw book updates [ID, price, qty], aggregated book updates [price, amount, count]
func NewBookUpdateFromRaw(symbol, precision string, data []interface{}) (b *BookUpdate, err error) {
	if len(data) < 3 {
		return b, fmt.Errorf("data slice too short for book update, expected %d got %d: %#v", 3, len(data), data)
	}
	var px decimal.Decimal
	var id, cnt int64
	amt := DecValOrZero(data[2])

	var side OrderSide
	var actionCtrl decimal.Decimal
	if IsRawBook(precision) {
		// [ID, price, amount]
		id = Int64ValOrZero(data[0])
		px = DecValOrZero(data[1])
		actionCtrl = px
	} else {
		// [price, amount, count]
		px = DecValOrZero(data[0])
		cnt = Int64ValOrZero(data[1])
		actionCtrl = decimal.New(cnt, 0)
	}

	if amt.Cmp(decimal.Zero) > 0 {
		side = Bid
	} else {
		side = Ask
	}

	var action BookAction
	if actionCtrl.Cmp(decimal.Zero) <= 0 {
		action = BookRemoveEntry
	} else {
		action = BookUpdateEntry
	}

	b = &BookUpdate{
		Symbol: symbol,
		Price:  px.Abs(),
		Count:  cnt,
		Amount: amt.Abs(),
		Side:   side,
		Action: action,
		ID:     id,
	}

	return
}

type Candle struct {
	Symbol     string
	Resolution CandleResolution
	MTS        int64
	Open       decimal.Decimal
	Close      decimal.Decimal
	High       decimal.Decimal
	Low        decimal.Decimal
	Volume     decimal.Decimal
}

type CandleSnapshot struct {
	Snapshot []*Candle
}

func toInterfaceSlice(nums []json.Number) []interface{} {
	data := make([]interface{}, len(nums))
	for j, f := range nums {
		data[j] = f
	}
	return data
}

func NewCandleSnapshotFromRaw(symbol string, resolution CandleResolution, raw [][]json.Number) (*CandleSnapshot, error) {
	if len(raw) <= 0 {
		return nil, fmt.Errorf("data slice too short for candle snapshot: %#v", raw)
	}
	snap := make([]*Candle, 0)
	for _, f := range raw {
		c, err := NewCandleFromRaw(symbol, resolution, toInterfaceSlice(f))
		if err == nil {
			snap = append(snap, c)
		}
	}
	return &CandleSnapshot{Snapshot: snap}, nil
}

func NewCandleFromRaw(symbol string, resolution CandleResolution, raw []interface{}) (c *Candle, err error) {
	if len(raw) < 6 {
		return c, fmt.Errorf("data slice too short for candle, expected %d got %d: %#v", 6, len(raw), raw)
	}

	c = &Candle{
		Symbol:     symbol,
		Resolution: resolution,
		MTS:        Int64ValOrZero(raw[0]),
		Open:       DecValOrZero(raw[1]),
		Close:      DecValOrZero(raw[2]),
		High:       DecValOrZero(raw[3]),
		Low:        DecValOrZero(raw[4]),
		Volume:     DecValOrZero(raw[5]),
	}

	return
}
