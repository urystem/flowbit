package domain

type Exchange struct {
	Source    string
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type ExchangeAggregation struct {
	Source   string
	Symbol   string
	Count    uint
	AvgPrice float64
	MinPrice float64
	MaxPrice float64
}
