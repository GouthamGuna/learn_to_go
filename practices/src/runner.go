package src

import (
	"fmt"
	"os"
)

func ApplicationRunner() {
	limits := 99
	fmt.Println("\nWelcome to the go programming language practice session")

	for limits < 100 {

		var options int

		fmt.Print("\nPress 1 : Investment Calculator \nPress 2 : Profit Calculator \nPress 3 : Quit \n\nEnter the values : ")

		fmt.Scan(&options)

		switch options {
		case (1):
			fmt.Println(" ")
			LogWritter("Invoking the InvestmentCalculator() fun...")
			InvestmentCalculator()
		case (2):
			fmt.Println(" ")
			LogWritter("Invoking the ProfitCalculator() fun...")
			ProfitCalculator()
		case (3):
			return
		default:
			LogWritter("Invoking the default block...")
			fmt.Println("Invalid Option")
			continue
		}
	}
}

func LogWritter(userLog string) {
	os.WriteFile("go_app.log", []byte(userLog), 0644)
}
