package src

import (
	"fmt"
	"os"
)

func ApplicationRunner() {
	limits := 99

	for limits < 100 {

		var options int

		fmt.Println("Welcome to the go programming language practice session")
		fmt.Print("\nPress 1 : Investment Calculator \nPress 2 : Profit Calculator \nPress 3 : Quit \nEnter the values : ")

		fmt.Scan(&options)

		switch options {
		case (1):
			fmt.Println(" ")
			InvestmentCalculator()
		case (2):
			fmt.Println(" ")
			ProfitCalculator()
		case (3):
			os.Exit(0)
		default:
			fmt.Println("Invalid Option")
		}
	}
}
