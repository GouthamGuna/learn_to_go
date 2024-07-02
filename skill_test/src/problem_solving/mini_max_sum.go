package problemsolving

import (
	"errors"
	"fmt"
	"sort"
)

func findMinMaxSum(arr []int64) ([2]int64, error) {
	if len(arr) < 4 {
		return [2]int64{}, errors.New("the input array must contain at least four elements")
	}

	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	minSum := int64(0)
	for i := 0; i < 4; i++ {
		minSum += arr[i]
	}

	maxSum := int64(0)
	for i := len(arr) - 4; i < len(arr); i++ {
		maxSum += arr[i]
	}

	return [2]int64{minSum, maxSum}, nil
}

func run_findMinMaxSum() {

	result, err := findMinMaxSum([]int64{1, 2, 3, 4, 5})

	if err != nil {
		fmt.Println("Error : ", err)
		return
	}

	fmt.Printf("\n Min Sum %d\n Max Sum %d\n", result[0], result[1])
}
