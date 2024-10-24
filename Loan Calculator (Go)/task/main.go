package main

import (
	"flag"
	"fmt"
	. "math"
	_ "os"
	_ "strings"
)

func main() {

	principal := flag.Float64("principal", 0.00, "Enter the principal amount")
	payment := flag.Float64("payment", 0.00, "Enter the payment amount")
	n := flag.Int("periods", 0, "Enter the time period")
	interestRate := flag.Float64("interest", 0.00, "Enter the interest amount")
	loanType := flag.String("type", "", "Specify the loan calculation type")

	flag.Parse()

	/*
		if !strings.HasPrefix(os.Args[1], "--type") {
			log.Fatal("Incorrect parameters")
		}
	*/

	switch *loanType {
	case "annuity":
		promptAnnuity(*principal, *interestRate, *payment, *n)
	case "diff":
		if *n < 1 || *principal < 0 || *interestRate <= 0 || *payment < 0 {
			fmt.Println("Incorrect parameters")
			break
		}
		var overPayment int

		for m := 1; m <= *n; m++ {
			diffPay := int(promptDifferentiated(*principal, *interestRate, *n, m))
			overPayment += diffPay
			fmt.Printf("Month %d: payment is %d\n", m, diffPay)
		}

		fmt.Printf("\nOverpayment = %d\n", overPayment-int(*principal))
	default:
		fmt.Println("Incorrect parameters")
	}

}

func calculateNumberOfPayments(payment, principal, interestRate float64) int {
	i := interestRate / (12 * 100) // Convert annual interest rate to monthly and to a decimal

	n := Log(payment/(payment-i*principal)) / Log(1+i)

	return int(Ceil(n)) // Round up to the next whole number
}

func calculateAnnuity(principal, interestRate float64, periods int) float64 {
	i := interestRate / (12 * 100)
	numerator := i * Pow(1.00+i, float64(periods))
	denominator := Pow(1.00+i, float64(periods)) - 1.00
	return Ceil(principal * (numerator / denominator))
}

func calculatePrincipal(payment, interestRate float64, periods int) float64 {
	i := interestRate / float64(12*100)
	numerator := i * Pow(1.00+i, float64(periods))
	denominator := Pow(1.00+i, float64(periods)) - 1.00

	return payment / (numerator / denominator)
}

func promptAnnuity(principal, interestRate, payment float64, periods int) {
	if principal < 0 || interestRate <= 0 || payment < 0 {
		fmt.Println("Incorrect parameters")
		return
	}

	switch {
	case payment == 0.00:
		payment = calculateAnnuity(principal, interestRate, periods)
		fmt.Printf("Your monthly payment = %d!\n", int(payment))
		fmt.Printf("Overpayment = %d\n", periods*int(payment)-int(principal))
	case periods == 0:
		periods = calculateNumberOfPayments(payment, principal, interestRate)
		if periods < 12 && periods > 0 {
			fmt.Printf("It will take %d months to repay this loan!\n", periods)
			fmt.Printf("Overpayment = %d\n", periods*int(payment)-int(principal))
		} else if periods%12 == 0 {
			fmt.Printf("It will take %d years to repay this loan!\n", periods/12)
			fmt.Printf("Overpayment = %d\n", periods*int(payment)-int(principal))
		} else {
			remainder := periods % 12
			fmt.Printf("It will take %d years and %d months to repay this loan!\n", periods/12, remainder)
			fmt.Printf("Overpayment = %d\n", periods*int(payment)-int(principal))
		}
	case principal == 0.00:
		principal = calculatePrincipal(payment, interestRate, periods)
		fmt.Printf("Your loan principal = %d!\n", int(principal))
		fmt.Printf("Overpayment = %d\n", periods*int(payment)-int(principal))
	default:
		fmt.Println("")
	}
}

func promptDifferentiated(principal, interestRate float64, periods, current int) float64 {
	i := interestRate / (12 * 100)
	bigFrac := principal - (float64(current-1)/float64(periods))*principal
	return Ceil((principal / float64(periods)) + i*bigFrac)
}
