package problemsolving

import (
	"errors"
)

func algorithmPlusMinus(arr []int64) ([3]float64, error) {

	if len(arr) <= 0 {
		return [3]float64{}, errors.New("the input array must contain at least one element")
	}

	var positive, negative, zero int64
	
	for _, num := range arr {
		if num > 0 {
			positive++
		} else if num < 0 {
			negative++
		} else {
			zero++
		}
	}

	totalElements := float64(len(arr))
	return [3]float64{
		float64(positive) / totalElements,
		float64(negative) / totalElements,
		float64(zero) / totalElements,
	}, nil
}
