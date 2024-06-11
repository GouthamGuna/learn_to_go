package src

import "fmt"

func ProfitCalculator() {
	var revenue float64
	var expenses float64
	var taxRate float64

	fmt.Print("Enter revenue :")
	fmt.Scan(&revenue)

	fmt.Print("Enter expenses :")
	fmt.Scan(&expenses)

	fmt.Print("Enter taxRate :")
	fmt.Scan(&taxRate)

	// earnings before tax
	ebt := revenue - expenses
	profit := ebt * (1 - taxRate/RatioBy)
	ratio := ebt / profit

	formattedEBT := fmt.Sprintf("Earnings before tax : %.2f", ebt)
	formattedProfit := fmt.Sprintf("\nProfit : %.2f", profit)
	formattedRatio := fmt.Sprintf("\nRatio : %.2f\n", ratio)

	fmt.Println(formattedEBT, formattedProfit, formattedRatio)
}
