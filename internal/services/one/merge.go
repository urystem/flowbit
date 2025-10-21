package one

import "marketflow/internal/domain"

func (one *oneMinute) merger(red, db []domain.ExchangeAggregation) []domain.ExchangeAggregation {
	for i, dv := range db {
		var found bool
		for j, rv := range red {
			if rv.Source == dv.Source && rv.Symbol == dv.Symbol {
				red[j].Count += dv.Count
				red[j].AvgPrice = (rv.AvgPrice*float64(rv.Count) + dv.AvgPrice*float64(dv.Count)) /
					float64(red[j].Count)
				red[j].MinPrice = min(rv.MinPrice, dv.MinPrice)
				red[j].MaxPrice = max(rv.MaxPrice, dv.MaxPrice)
				found = true
				break
			}
		}
		if !found {
			red = append(red, db[i])
		}
	}
	return red
}

// avgtotal​=​(avg1​×count1​+avg2​×count2)/(count1​+count2)​​
