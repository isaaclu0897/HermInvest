package main

import (
	"HermInvest/pkg/model"
	"HermInvest/pkg/service"
	"fmt"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   `query {--all | --id <ID> | [--stockNo <StockNumber> --type <Type> --date <Date>]}`,
	Short: "Query stock (Transaction ID, Stock No., Type, or Date)",
	Example: "" +
		"  - Query by Transaction ID:\n" +
		"    hermInvestCli stock query --id 11\n\n" +

		"  - Query all records:\n" +
		"    hermInvestCli stock query --all\n\n" +

		"  - Query by stock number:\n" +
		"    hermInvestCli stock query --stockNo 0050\n\n" +

		"  - Query by stock number, type, and date:\n" +
		"    hermInvestCli stock query --stockNo 0050 --type 1 --date 2023-12-01",
	Long: "Query stock by transaction ID, stock number, type, or date.",
	Args: cobra.NoArgs,
	RunE: queryRun,
}

func init() {
	stockCmd.AddCommand(queryCmd)

	queryCmd.Flags().Bool("all", false, "Query all records")
	queryCmd.Flags().Int("id", 0, "Query by ID")
	queryCmd.Flags().String("stockNo", "", "Stock number")
	queryCmd.Flags().Int("type", 0, "Type")
	queryCmd.Flags().String("date", "", "Date")
}

func queryRun(cmd *cobra.Command, args []string) error {
	if cmd.Flags().NFlag() == 0 {
		return fmt.Errorf("no flags provided")
	}

	all, _ := cmd.Flags().GetBool("all")
	id, _ := cmd.Flags().GetInt("id")
	stockNo, _ := cmd.Flags().GetString("stockNo")
	tranType, _ := cmd.Flags().GetInt("type")
	date, _ := cmd.Flags().GetString("date")

	serv := service.InitializeService()

	var transactions []*model.Transaction
	var transactionsErr error
	if all {
		transactions, transactionsErr = serv.QueryTransactionAll()
	} else if id != 0 {
		var transaction *model.Transaction
		transaction, transactionsErr = serv.QueryTransactionByID(id)
		transactions = append(transactions, transaction)
	} else {
		transactions, transactionsErr = serv.QueryTransactionByDetails(stockNo, tranType, date)
	}
	if transactionsErr != nil {
		fmt.Println("Error querying database:", transactionsErr)
	} else {
		displayResults(transactions)
	}

	return nil
}

func displayResults(transactions []*model.Transaction) {
	fmt.Print("ID,\tStock No,\tType,\tQty(shares),\tUnit Price,\tTotal Amount,\ttaxes\n")
	for _, t := range transactions {
		fmt.Printf("%d,\t%8s,\t%4d,\t%11d,\t%10.2f,\t%12d,\t%5d\n", t.ID, t.StockNo, t.TranType, t.Quantity, t.UnitPrice, t.TotalAmount, t.Taxes)
	}
}
