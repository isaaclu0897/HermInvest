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
	// 1. new transaction from input
	// 2. find the first purchase from the inventory
	// 3. check the transaction type of new transaction and first purchase

	// TODO: service.addTransaction() AddTransactionAndUpdateInventory
	newTransaction := model.NewTransactionFromInput(stockNo, date, quantity, tranType, unitPrice)

	t, err := repo.AddTransaction(newTransaction)
	if err != nil {
		fmt.Println("Error adding transaction: ", err)
	} else if t != nil {
		var ts []*model.Transaction
		ts = append(ts, t)
		displayResults(ts)
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
