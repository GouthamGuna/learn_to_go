package src

import "fmt"

func ProfitCalculator() {
	revenue := fetchResourceFromUser("Enter revenue : ")
	expenses := fetchResourceFromUser("Enter expenses : ")
	taxRate := fetchResourceFromUser("Enter taxRate : ")

	// earnings before tax
	ebt := revenue - expenses
	profit := ebt * (1 - taxRate/RatioBy)
	ratio := ebt / profit

	formattedEBT := fmt.Sprintf("Earnings before tax : %.2f", ebt)
	formattedProfit := fmt.Sprintf("\nProfit : %.2f", profit)
	formattedRatio := fmt.Sprintf("\nRatio : %.2f\n", ratio)

	fmt.Println(formattedEBT, formattedProfit, formattedRatio)
}
