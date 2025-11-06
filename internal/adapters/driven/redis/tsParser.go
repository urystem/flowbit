package redis

import "fmt"

func (r *myRedis) timeValParser(zat any) (int64, float64, error) {
	timeVal, ok := zat.([]any)
	if !ok {
		return 0, 0, fmt.Errorf("%s", "unkown error 1")
	} else if len(timeVal) != 2 {
		return 0, 0, fmt.Errorf("%s%v", "unkown error 2", timeVal)
	}
	myTime, ok := timeVal[0].(int64)
	if !ok {
		return 0, 0, fmt.Errorf("%s", "unkown error 3")
	}
	price, ok := timeVal[1].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("%s", "unkown error 4")
	}
	return myTime, price, nil
}

func (r *myRedis) getExchangeName(zat any) (string, error) {
	mapa, ok := zat.(map[any]any)
	if !ok {
		return "", fmt.Errorf("unexpected type (get exchange name)")
	}
	exAny, ok := mapa["exchange"]
	if !ok {
		return "", fmt.Errorf("not found exchange label")
	}
	ex, ok := exAny.(string)
	if !ok {
		return "", fmt.Errorf("unexpected type (parse to string)")
	}
	return ex, nil
}
