package domain

type Exchange struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
}
