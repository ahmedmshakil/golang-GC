package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)
type Transaction struct {
	Date        time.Time
	Type        string
	Category    string
	Amount      float64
	Description string
}

type Data struct {
	Transactions []Transaction
}
const (
	Income  = "Income"
	Expense = "Expense"
	Month   = "month"
	Year    = "year"
	All     = "all"
)
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
func parseFloat(amountStr string) (float64, error) {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %w", err)
	}
	return amount, nil
}

// add a new transaction
func (d *Data) addTransaction(date time.Time, transactionType, category string, amount float64, description string) error {
	if transactionType != Income && transactionType != Expense {
		return fmt.Errorf("invalid transaction type: %s", transactionType)
	}
	d.Transactions = append(d.Transactions, Transaction{Date: date, Type: transactionType, Category: category, Amount: amount, Description: description})
	return nil
}

func (d *Data) importTransactions(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV data: %w", err)
	}

	if len(records) <= 1 {
		return fmt.Errorf("empty or invalid CSV file")
	}

	for i, record := range records[1:] {
		if len(record) != 5 {
			fmt.Printf("Skipping record %d due to invalid number of fields: %v\n", i+2, record) 
			continue
		}
		date, err := parseDate(record[0])
		if err != nil {
			fmt.Printf("Skipping record %d due to invalid date: %v, error: %v\n", i+2, record, err)
			continue
		}
		transactionType := record[1]
		category := record[2]
		amount, err := parseFloat(record[3])
		if err != nil {
			fmt.Printf("Skipping record %d due to invalid amount: %v, error: %v\n", i+2, record, err)
			continue
		}
		description := record[4]

		err = d.addTransaction(date, transactionType, category, amount, description)
		if err != nil {
			fmt.Printf("Skipping record %d due to error: %v, error: %v \n", i+2, record, err)
			continue
		}
	}
	return nil
}
func (d *Data) calculateSummary(period string, periodValue string) (float64, float64, map[string]float64) {
	totalIncome := 0.0
	totalExpenses := 0.0
	categorySummary := make(map[string]float64)

	for _, transaction := range d.Transactions {
		include := false
		switch period {
		case Month:
			inputTime, _ := time.Parse("2006-01", periodValue)
			if transaction.Date.Year() == inputTime.Year() && transaction.Date.Month() == inputTime.Month() {
				include = true
			}
		case Year:
			inputTime, _ := time.Parse("2006", periodValue)
			if transaction.Date.Year() == inputTime.Year() {
				include = true
			}
		case All:
			include = true
		}

		if include {
			if transaction.Type == Income {
				totalIncome += transaction.Amount
			} else if transaction.Type == Expense {
				totalExpenses += transaction.Amount
			}
			categorySummary[transaction.Category] += transaction.Amount
		}
	}
	return totalIncome, totalExpenses, categorySummary
}
func (d *Data) displaySummary(period string, periodValue string) {
	totalIncome, totalExpenses, categorySummary := d.calculateSummary(period, periodValue)
	fmt.Printf("Income: %.2f\n", totalIncome)
	fmt.Printf("Expenses: %.2f\n", totalExpenses)
	fmt.Printf("Net Balance: %.2f\n", totalIncome-totalExpenses)
	fmt.Println("Category Summary:")
	for category, amount := range categorySummary {
		fmt.Printf("  %s: %.2f\n", category, amount)
	}
}
func (d *Data) predictExpenses(months int) ([]float64, []float64) {
	expenses := make([]float64, 0)
	for _, transaction := range d.Transactions {
		if transaction.Type == Expense {
			expenses = append(expenses, transaction.Amount)
		}
	}

	predictedExpenses := make([]float64, months)
	predictedNetBalance := make([]float64, months)
    totalIncome, _, _ := d.calculateSummary(All, "") 
	if len(expenses) > 0 {
		lastExpense := expenses[len(expenses)-1]
		for i := 0; i < months; i++ {
			predictedExpenses[i] = lastExpense * (1 + 0.1*float64(i+1)) // 10% increase per month
            predictedNetBalance[i] = totalIncome - predictedExpenses[i]
		}
	} else {
        for i := 0; i < months; i++{
            predictedExpenses[i] = 0
            predictedNetBalance[i] = totalIncome
        }
    }
	return predictedExpenses, predictedNetBalance
}
func (d *Data) displayPredictions(months int) {
	predictedExpenses, predictedNetBalance := d.predictExpenses(months)
	fmt.Println("Predicted Expenses for the next", months, "months:")
	for i, expense := range predictedExpenses {
		fmt.Printf("  Month %d: %.2f\n", i+1, expense)
	}
    fmt.Println("Predicted Net Balance for the next", months, "months:")
	for i, balance := range predictedNetBalance{
		fmt.Printf("  Month %d: %.2f\n", i+1, balance)
	}
}

//display
func displayHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  add    Add a new transaction")
	fmt.Println("  import Import transactions from a CSV file")
	fmt.Println("  summary Display a summary of income, expenses, and net balance")
	fmt.Println("  predict Display predicted expenses and net balance")
	fmt.Println("  help   Display this help message")
	fmt.Println("  exit   Exit the application")
}

func main() {
	data := Data{}
	fmt.Println("Welcome to Personal Finance Tracker!")
	displayHelp()

	for {
		fmt.Print("\nEnter command: ")
		var command string
		fmt.Scanln(&command)

		command = strings.ToLower(command) 

		switch command {
		case "add":
			var dateStr, transactionType, category, description string
			var amountStr string

			fmt.Print("Date (YYYY-MM-DD): ")
			fmt.Scanln(&dateStr)
			date, err := parseDate(dateStr)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			fmt.Print("Type (Income/Expense): ")
			fmt.Scanln(&transactionType)

			fmt.Print("Category: ")
			fmt.Scanln(&category)

			fmt.Print("Amount: ")
			fmt.Scanln(&amountStr)
			amount, err := parseFloat(amountStr)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			fmt.Print("Description: ")
			fmt.Scanln(&description)

			err = data.addTransaction(date, transactionType, category, amount, description)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Transaction added successfully.")
			}

		case "import":
			var filename string
			fmt.Print("Enter CSV filename: ")
			fmt.Scanln(&filename)
			err := data.importTransactions(filename)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Transactions imported successfully.")
			}

		case "summary":
			var period, periodValue string
			fmt.Print("Time period (month/year/all): ")
			fmt.Scanln(&period)
			period = strings.ToLower(period) //forgiving input

			switch period {
			case Month:
				fmt.Print("Month (YYYY-MM): ")
				fmt.Scanln(&periodValue)
				if _, err := time.Parse("2006-01", periodValue); err != nil {
					fmt.Println("Error: Invalid month format. Please use YYYY-MM.")
					break
				}
			case Year:
				fmt.Print("Year (YYYY): ")
				fmt.Scanln(&periodValue)
				if _, err := time.Parse("2006", periodValue); err != nil {
					fmt.Println("Error: Invalid year format. Please use YYYY.")
					break
				}
			case All:
				periodValue = ""
			default:
				fmt.Println("Error: Invalid time period. Please use month, year, or all.")
				break
			}
			data.displaySummary(period, periodValue)

		case "predict":
			var months int
			fmt.Print("Prediction period (months): ")
			fmt.Scanln(&months)
			if months <= 0 {
				fmt.Println("Error: Number of months must be greater than zero.")
				break
			}
			data.displayPredictions(months)

		case "help":
			displayHelp()

		case "exit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid command. Please try again.")
			displayHelp()
		}
	}
}

