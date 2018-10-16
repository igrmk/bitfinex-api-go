package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestTickerAll(t *testing.T) {
	httpDo := func(_ *http.Client, req *http.Request) (*http.Response, error) {
		msg := `[["tSYMBOL1",0.01,0.02,0.03,0.04,0.05,0.06,0.07,0.08,0.09,0.10],["tSYMBOL2",0.11,0.12,0.13,0.14,0.15,0.16,0.17,0.18,0.19,0.50]]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	ticker, err := NewClientWithHttpDo(httpDo).Ticker.All()

	if err != nil {
		t.Fatal(err)
	}

	if len(ticker.Snapshot) != 2 {
		t.Fatalf("expected 2 ticker entries, but got %d", len(ticker.Snapshot))
	}

	if ticker.Snapshot[1].Symbol != "symbol2" {
		t.Fatalf("expected symbol2 symbol, but got %s", ticker.Snapshot[1].Symbol)
	}

	if ticker.Snapshot[1].Low != 0.5 {
		t.Fatalf("expected low equal to 0.5, but got %f", ticker.Snapshot[1].Low)
	}
}
