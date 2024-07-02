package problemsolving

import "fmt"

func PS_PackageRunner() {

	result, err := algorithmPlusMinus([]int64{1, 1, 0, -1, -1})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf(" Positive fraction: %.6f\n Negative fraction: %.6f\n Zero fraction: %.6f\n ", result[0], result[1], result[2])
}
