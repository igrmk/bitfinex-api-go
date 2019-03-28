package bitfinex

import (
	"encoding/json"
	"fmt"

	"github.com/igrmk/decimal"
)

func JsonNumberSlice(in []interface{}) ([]json.Number, error) {
	var ret []json.Number
	for _, e := range in {
		if item, ok := e.(json.Number); ok {
			ret = append(ret, item)
		} else {
			return nil, fmt.Errorf("expected slice of numbers but got: %v", in)
		}
	}

	return ret, nil
}

func DecValOrZero(i interface{}) decimal.Decimal {
	if num, ok := i.(json.Number); ok {
		return decimal.RequireFromString(num.String())
	}
	return decimal.Zero
}

func Int64ValOrZero(i interface{}) int64 {
	if num, ok := i.(json.Number); ok {
		if parsed, err := num.Int64(); err == nil {
			return parsed
		}
	}
	return 0
}

func bValOrFalse(i interface{}) bool {
	if r, ok := i.(bool); ok {
		return r
	}
	return false
}

func strValOrEmpty(i interface{}) string {
	if r, ok := i.(string); ok {
		return r
	}
	return ""
}
