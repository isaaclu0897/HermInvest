package main

import (
	"HermInvest/pkg/repository"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use: "web",
	Run: webRun,
}

func init() {
	stockCmd.AddCommand(webCmd)
}

func webRun(cmd *cobra.Command, args []string) {
	router := gin.Default()

	// load HTML files - Is it a best way?
	router.LoadHTMLFiles("html/transactionDetails.html")

	router.GET("/", homePage)
	router.GET("/transaction", transactionPage)
	router.GET("/api/transaction", apiGetTransactions)
	router.GET("/transactionDetails/:stockNo", transactionDetailsPage)
	router.GET("/api/transaction/:stockNo", apiGetTransactionsByStockNo)
	router.Static("/assets", "./assets")

	open("http://127.0.0.1:9453/transaction")

	err := router.Run(":9453")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func homePage(c *gin.Context) {
	var pageHTML []byte
	pageHTML, err := os.ReadFile("html/home.html")
	if err != nil {
		log.Fatal("os.ReadFile: ", err)
	}

	c.Data(http.StatusOK, "text/html", pageHTML)
}

func apiGetTransactions(c *gin.Context) {
	// transactions := []*model.Transaction{
	// 	{
	// 		ID:          0,
	// 		Date:        "",
	// 		Time:        "",
	// 		StockNo:     "0050",
	// 		TranType:    0,
	// 		Quantity:    5000,
	// 		UnitPrice:   111,
	// 		TotalAmount: 557550,
	// 		Taxes:       1672,
	// 		StockMapping: model.StockMapping{
	// 			StockNo:   "0050",
	// 			StockName: "元大台灣50",
	// 		},
	// 	},
	// 	{
	// 		ID:          0,
	// 		Date:        "",
	// 		Time:        "",
	// 		StockNo:     "00902",
	// 		TranType:    0,
	// 		Quantity:    7000,
	// 		UnitPrice:   12,
	// 		TotalAmount: 87010,
	// 		Taxes:       259,
	// 		StockMapping: model.StockMapping{
	// 			StockNo:   "00902",
	// 			StockName: "中信電池及儲能",
	// 		},
	// 	},
	// }
	db := repository.GetDBConnection()

	// init transactionRepository
	repo := repository.NewRepository(db)

	transactions, err := repo.QueryTransactionInventory()
	if err != nil {
		fmt.Println("err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query transaction"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func apiGetTransactionsByStockNo(c *gin.Context) {
	db := repository.GetDBConnection()

	repo := repository.NewRepository(db)

	stockNo := c.Param("stockNo")

	transactions, err := repo.QueryTransactionInventoryByStockNo(stockNo)
	if err != nil {
		fmt.Println("err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query transaction by StockNo"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func transactionPage(c *gin.Context) {

	var pageHTML []byte
	pageHTML, err := os.ReadFile("html/transaction.html")
	if err != nil {
		log.Fatal("os.ReadFile: ", err)
	}

	c.Data(http.StatusOK, "text/html", pageHTML)
}

func transactionDetailsPage(c *gin.Context) {

	stockNo := c.Param("stockNo")

	// transfer stockNo to template
	c.HTML(http.StatusOK, "transactionDetails.html", gin.H{
		"stockNo": stockNo,
	})
}

func open(url string) error { // open url from browser
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
