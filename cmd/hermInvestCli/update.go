package main

import (
	"HermInvest/pkg/service"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// 1. check input
// 2. query transacation by id
// 3. recalc total amount and taxes
// 4. build sql syntax
// 5. insert into sql
// 6. print out result

var updateCmd = &cobra.Command{
	Use:   "update id unitPrice",
	Short: "Update unit price by transaction ID",
	Example: "" +
		"  - Update unit Price by ID:\n" +
		"    hermInvestCli stock update 11 20.3",
	Long: `Update the unit price of stock in the inventory using the transaction ID.`,
	Args: cobra.ExactArgs(2),
	Run:  updateRun,
}

func init() {
	stockCmd.AddCommand(updateCmd)
}

func updateRun(cmd *cobra.Command, args []string) {
	transactionID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error parsing integer: ", err)
	}
	unitPrice, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Println("Error parsing float: ", err)
	}

	serv := service.InitializeService()

	t, err := serv.QueryTransactionByID(transactionID)
	if err != nil {
		fmt.Println("Error querying database:", err)
	}

	// TODO: check update work. ex: update a fake transaction ID to db
	t.SetUnitPrice(unitPrice) // update unit Price

	// update db
	err = serv.UpdateTransaction(t.ID, t)
	if err != nil {
		fmt.Println("Error updating stock information:", err)
		return
	}

	fmt.Printf("Successfully updated transaction ID %d with new unit price %.2f\n", t.ID, t.UnitPrice)
}
