package domain

import "time"

type Exchange struct {
	Source    string
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type ExchangeAggregation struct {
	Source   string
	Symbol   string
	Count    int
	AvgPrice float64
	MinPrice float64
	MaxPrice float64
}

type ExchangeAvg struct {
	ExAvgs []ExchangeAggregation
	AtTime time.Time
}
