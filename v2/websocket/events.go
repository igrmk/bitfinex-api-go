package websocket

import (
	"encoding/json"
	"fmt"
)

type eventType struct {
	Event string `json:"event"`
}

type InfoEvent struct {
	Version float64 `json:"version"`
}

type RawEvent struct {
	Data interface{}
}

type AuthEvent struct {
	Event   string       `json:"event"`
	Status  string       `json:"status"`
	ChanID  int64        `json:"chanId,omitempty"`
	UserID  int64        `json:"userId,omitempty"`
	SubID   string       `json:"subId"`
	AuthID  string       `json:"auth_id,omitempty"`
	Message string       `json:"msg,omitempty"`
	Caps    Capabilities `json:"caps"`
}

type Capability struct {
	Read  int `json:"read"`
	Write int `json:"write"`
}

type Capabilities struct {
	Orders    Capability `json:"orders"`
	Account   Capability `json:"account"`
	Funding   Capability `json:"funding"`
	History   Capability `json:"history"`
	Wallets   Capability `json:"wallets"`
	Withdraw  Capability `json:"withdraw"`
	Positions Capability `json:"positions"`
}

// error codes pulled from v2 docs & API usage
const (
	ErrorCodeUnknownEvent         int = 10000
	ErrorCodeUnknownPair          int = 10001
	ErrorCodeUnknownBookPrecision int = 10011
	ErrorCodeUnknownBookLength    int = 10012
	ErrorCodeSubscriptionFailed   int = 10300
	ErrorCodeAlreadySubscribed    int = 10301
	ErrorCodeUnknownChannel       int = 10302
	ErrorCodeUnsubscribeFailed    int = 10400
	ErrorCodeNotSubscribed        int = 10401
)

type ErrorEvent struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`

	// also contain members related to subscription reject
	SubID     string `json:"subId"`
	Channel   string `json:"channel"`
	ChanID    int64  `json:"chanId"`
	Symbol    string `json:"symbol"`
	Precision string `json:"prec,omitempty"`
	Frequency string `json:"freq,omitempty"`
	Key       string `json:"key,omitempty"`
	Len       string `json:"len,omitempty"`
	Pair      string `json:"pair"`
}

type UnsubscribeEvent struct {
	Status string `json:"status"`
	ChanID int64  `json:"chanId"`
}

type SubscribeEvent struct {
	SubID     string `json:"subId"`
	Channel   string `json:"channel"`
	ChanID    int64  `json:"chanId"`
	Symbol    string `json:"symbol"`
	Precision string `json:"prec,omitempty"`
	Frequency string `json:"freq,omitempty"`
	Key       string `json:"key,omitempty"`
	Len       string `json:"len,omitempty"`
	Pair      string `json:"pair"`
}

type ConfEvent struct {
	Flags int `json:"flags"`
}

// onEvent handles all the event messages and connects SubID and ChannelID.
func (c *Client) handleEvent(msg []byte) error {
	basicEvent := &eventType{}
	err := json.Unmarshal(msg, basicEvent)
	if err != nil {
		return err
	}
	//var e interface{}
	switch basicEvent.Event {
	case "info":
		ev := InfoEvent{}
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			return err
		}
		err_open := c.handleOpen()
		if err_open != nil {
			return err_open
		}
		c.listener <- &ev
	case "auth":
		ev := AuthEvent{}
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			return err
		}
		if ev.Status != "" && ev.Status == "OK" {
			c.Authentication = SuccessfulAuthentication
		} else {
			c.Authentication = RejectedAuthentication
		}
		c.handleAuthAck(&ev)
		c.listener <- &ev
		return nil
	case "subscribed":
		ev := SubscribeEvent{}
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			return err
		}
		err = c.subscriptions.activate(ev.SubID, ev.ChanID)
		if err != nil {
			return err
		}
		c.listener <- &ev
		return nil
	case "unsubscribed":
		ev := UnsubscribeEvent{}
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			return err
		}
		err_rem := c.subscriptions.removeByChannelID(ev.ChanID)
		if err_rem != nil {
			return err_rem
		}
		c.listener <- &ev
	case "error":
		ev := ErrorEvent{}
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			return err
		}
		c.listener <- &ev
	case "conf":
		ev := ConfEvent{}
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			return err
		}
		c.listener <- &ev
	default:
		return fmt.Errorf("unknown event: %s", msg) // TODO: or just log?
	}

	//err = json.Unmarshal(msg, &e)
	//TODO raw message isn't ever published

	return err
}
