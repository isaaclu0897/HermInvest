package main

import (
	"HermInvest/pkg/model"
	"HermInvest/pkg/repository"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// 1. check input
// 2. calc total amount and taxes
// 3. build sql syntax
// 4. insert into sql
// 5. print out result

var addCmd = &cobra.Command{
	Use:   "add stockNo type quantity unitPrice [date]",
	Short: "Add stock (Stock No., Type, Quantity, Unit Price)",
	Example: "" +
		"  - Purchase at today's date:\n" +
		"    hermInvestCli stock add 0050 1 1500 23.5\n\n" +

		"  - Sale on a specific date:\n" +
		"    hermInvestCli stock add -- 0050 -1 1500 23.5 2023-12-01",
	Long: `Add stock by transaction stock number, type, quantity, and unit price`,
	Args: cobra.RangeArgs(4, 5),
	Run:  addRun,
}

func init() {
	stockCmd.AddCommand(addCmd)
}

func addRun(cmd *cobra.Command, args []string) {
	stockNo, tranType, quantity, unitPrice, date, err := ParseTransactionForAddCmd(args)
	if err != nil {
		fmt.Println("Error parsing transaction data:", err)
		return
	}

	db, err := repository.GetDBConnection()
	if err != nil {
		fmt.Println("Error geting DB connection: ", err)
	}
	defer db.Close()

	// init transactionRepository
	repo := repository.NewTransactionRepository(db)

	// add stock in inventory
	t := model.NewTransactionFromInput(stockNo, date, quantity, tranType, unitPrice)

	if t.TranType < 0 { // sale stock
		firstPurchase, err := repo.FindFirstPurchase(stockNo)
		if err != nil {
			// handle find first purchase failed
		}

		// update stock quantity
		updateQuantity := firstPurchase.Quantity - t.Quantity
		firstPurchase.SetQuantity(updateQuantity)
		err = repo.UpdateTransaction(firstPurchase.ID, firstPurchase)
		if err != nil {
			// handle update transaction failed
		}

		// record purchase history
		purchaseHistory := firstPurchase
		purchaseHistory.SetQuantity(t.Quantity)
		_, err = repo.CreateTransactionHistory(purchaseHistory)
		if err != nil {
			// handle create transaction history failed
		}

		// record sale history
		saleHistory := t
		_, err = repo.CreateTransactionHistory(saleHistory)
		if err != nil {
			// handle create transaction history failed
		}

	} else {
		id, err := repo.CreateTransaction(t)
		if err != nil {
			fmt.Println("Error creating transaction: ", err)
		}
		transactions, err := repo.QueryTransactionByID(id)
		if err != nil {
			fmt.Println("Error querying database:", err)
		}

		// Print out result
		displayResults(transactions)
	}
}

func ParseTransactionForAddCmd(args []string) (string, int, int, float64, string, error) {
	stockNo := args[0] // regex a-z 0-9

	tranType, err := strconv.Atoi(args[1])
	if err != nil {
		return "", 0, 0, 0, "", fmt.Errorf("error parsing integer: %s", err)
	}

	quantity, err := strconv.Atoi(args[2])
	if err != nil {
		return "", 0, 0, 0, "", fmt.Errorf("error parsing integer: %s", err)
	}

	unitPrice, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return "", 0, 0, 0, "", fmt.Errorf("error parsing float: %s", err)
	}

	var date string
	if len(args) > 4 {
		parsedTime, err := time.Parse(time.DateOnly, args[4])
		if err != nil {
			return "", 0, 0, 0, "", fmt.Errorf("error parsing date: %s", err)
		}
		date = parsedTime.Format(time.DateOnly)
	} else {
		date = time.Now().Format(time.DateOnly)
	}

	return stockNo, tranType, quantity, unitPrice, date, nil
}
