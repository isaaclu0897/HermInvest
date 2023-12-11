package main

import (
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
	Use:   "add id stockNo type quantity unitPrice [date]",
	Short: "Add stock",
	Long:  `Add stock to the inventory`,
	Args:  cobra.MinimumNArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error parsing integer: ", err)
		}
		stockNo := args[1]
		tranType, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Error parsing integer: ", err)
		}
		quantity, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("Error parsing integer: ", err)
		}
		unitPrice, err := strconv.ParseFloat(args[4], 64)
		if err != nil {
			fmt.Println("Error parsing float: ", err)
		}

		// Parse date argument or default is today's date
		var date string
		if len(args) > 5 {
			date = args[5] // check user input time format is correct
		} else {
			date = time.Now().Format("2006-01-02")
		}

		t := NewTransactionFromUserInput(id, stockNo, date, quantity, tranType, unitPrice)

		db, err := GetDBConnection()
		if err != nil {
			fmt.Println("Error geting DB connection: ", err)
		}
		defer db.Close()

		// Execute the insert query
		query := `INSERT INTO tblTransaction (id, stockNo, date, quantity, tranType, unitPrice, totalAmount, taxes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = db.Exec(query, t.id, t.stockNo, t.date, t.quantity, t.tranType, t.unitPrice, t.totalAmount, t.taxes)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			fmt.Println("Pass: Stock added successfully!")
		}

		// Print out result
		rows, err := db.Query(buildQueryByID(id))
		if err != nil {
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		displayResults(rows)
	},
}
