package redis

// import (
// 	"context"
// 	"fmt"
// 	"strings"

// 	"marketflow/internal/domain"
// )

// func (rdb *myRedis) GetByLabel(ctx context.Context, from, to int, keys ...string) ([]domain.Exchange, error) {
// 	res, err := rdb.TSMRange(
// 		ctx,
// 		from,
// 		to,
// 		keys,
// 	).Result()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var exchanges []domain.Exchange
// 	// serie это 3 интерфейс [3]interface
// 	// и нам нужна только последный
// 	for exSym, serie := range res {
// 		points, ok := serie[2].([]any)
// 		if !ok {
// 			return nil, fmt.Errorf("%s", "series is not [3]any")
// 		}
// 		myexSym := strings.Split(exSym, ":")
// 		if len(myexSym) != 2 {
// 			return nil, fmt.Errorf("%s", "exchange and symbol key not valid")
// 		}
// 		for _, p := range points {
// 			arr, ok := p.([]any)
// 			if !ok || len(arr) != 2 {
// 				return nil, fmt.Errorf("%s", "here is not time and price")
// 			}
// 			myTime, ok := arr[0].(int64)
// 			if !ok {
// 				return nil, fmt.Errorf("%s", "invalid time")
// 			}
// 			myPrice, ok := arr[1].(float64)
// 			if !ok {
// 				return nil, fmt.Errorf("%s", "invalid price")
// 			}
// 			exchanges = append(exchanges, domain.Exchange{
// 				Source:    myexSym[0],
// 				Symbol:    myexSym[1],
// 				Price:     myPrice,
// 				Timestamp: myTime,
// 			})
// 		}
// 	}
// 	return exchanges, nil
// }
