package src

import (
	"fmt"
	"math"
)

const RatioBy float64 = 100.0

func DoCalculation() {

	const inflationRate = 6.5

	var investmentAmount float64
	var years float64

	expectedReturnRate := 5.5

	fmt.Print("Enter Investing Amount : ")
	fmt.Scan(&investmentAmount)

	fmt.Print("Enter Year : ")
	fmt.Scan(&years)

	fmt.Print("Enter Expected Return Rate : ")
	fmt.Scan(&expectedReturnRate)

	futureValue := investmentAmount * math.Pow(1+expectedReturnRate/RatioBy, years)
	futureRealValue := futureValue / math.Pow(1+inflationRate/100, years)

	fmt.Printf("futureValue : %.2f, futureRealValue : %.2f", futureValue, futureRealValue)
}
