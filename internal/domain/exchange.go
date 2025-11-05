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
	Source    string    `json:"source,omitempty"`
	Symbol    string    `json:"symbol,omitempty"`
	Count     uint      `json:"count,omitempty"`
	AvgPrice  float64   `json:"average_price,omitempty"`
	MinPrice  float64   `json:"min_price,omitempty"`
	MaxPrice  float64   `json:"max_price,omitempty"`
	Timestamp time.Time `json:"timestamp,omitzero"`
}
