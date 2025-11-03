package domain

import "time"

type Exchange struct {
	Source    string  `json:"source"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

type GetExchange struct {
	*Exchange
	Timestamp time.Time `json:"timestamp,omitzero"`
	Info      string    `json:"info,omitempty"`
}

type ExchangeAggregation struct {
	Source   string
	Symbol   string
	Count    uint
	AvgPrice float64
	MinPrice float64
	MaxPrice float64
}
