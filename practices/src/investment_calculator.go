package src

import (
	"fmt"
	"math"
)

const RatioBy float64 = 100.0
const inflationRate = 6.5

func InvestmentCalculator() {
	expectedReturnRate := 5.5

	investmentAmount := fetchResourceFromUser("Enter Investing Amount : ")
	years := fetchResourceFromUser("Enter Year : ")
	expectedReturnRate = fetchResourceFromUser("Enter Expected Return Rate : ")

	futureValue, futureRealValue := doCalculate(investmentAmount, expectedReturnRate, years)

	formattedFV := fmt.Sprintf("Future Value: %.2f", futureValue)
	formattedRFV := fmt.Sprintf("\nFuture Value (adjusted for Inflation): %.2f\n", futureRealValue)

	fmt.Println(formattedFV, formattedRFV)
}

func doCalculate(investmentAmount, expectedReturnRate, years float64) (float64, float64) {
	futureValue := investmentAmount * math.Pow(1+expectedReturnRate/RatioBy, years)
	futureRealValue := futureValue / math.Pow(1+inflationRate/100, years)
	return futureValue, futureRealValue
}

func fetchResourceFromUser(infoText string) (responseValue float64) {
	fmt.Print(infoText)
	fmt.Scan(&responseValue)
	return responseValue
}
