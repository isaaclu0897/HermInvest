package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Create SQLite Database
	db, err := sql.Open("sqlite3", "./internal/app/database/dev-database.db")
	if err != nil {
		fmt.Println("Error creating database:", err)
		return
	}
	defer db.Close()
	fmt.Println("Database is creating ...")

	// Create tblTransactionRecord table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "tblTransactionRecord" (
			"date"	TEXT NOT NULL,
			"time"	TEXT NOT NULL,
			"stockNo"	TEXT NOT NULL,
			"stockName"	TEXT NOT NULL,
			"tranType"	INTEGER NOT NULL,
			"quantity"	INTEGER NOT NULL,
			"unitPrice"	REAL NOT NULL,
			"source"	INTEGER NOT NULL,
			PRIMARY KEY("date","time")
		)
	`)
	if err != nil {
		fmt.Println("Error creating tblTransactionRecord table:", err)
		return
	}
	fmt.Println("Table tblTransactionRecord created successfully")

	// Create tblTransactionRecordSys table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "tblTransactionRecordSys" (
			"date"	TEXT NOT NULL,
			"time"	TEXT NOT NULL,
			"stockNo"	TEXT NOT NULL,
			"tranType"	INTEGER NOT NULL,
			"quantity"	INTEGER NOT NULL,
			"unitPrice"	REAL NOT NULL
		)
	`)
	if err != nil {
		fmt.Println("Error creating tblTransactionRecordSys table:", err)
		return
	}
	fmt.Println("Table tblTransactionRecordSys created successfully")

	// Create tblCapitalReduction table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS  tblCapitalReduction (
			YQ TEXT NOT NULL,
			stockNo TEXT NOT NULL,
			capitalReductionDate TEXT NOT NULL,
			distributionDate TEXT NOT NULL,
			cash REAL,
			ratio REAL,
			newStockNo TEXT
		)
	`)
	if err != nil {
		fmt.Println("Error creating tblCapitalReduction table:", err)
		return
	}
	fmt.Println("Table tblCapitalReduction created successfully")

	// Create tblDividend table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tblDividend (
			YQ TEXT NOT NULL,
			stockNo TEXT NOT NULL,
			ExDividendDate TEXT NOT NULL,
			distributionDate TEXT NOT NULL,
			cashDividend REAL,
			stockDividend REAL
		)
	`)
	if err != nil {
		fmt.Println("Error creating tblDividend table:", err)
		return
	}
	fmt.Println("Table tblDividend created successfully")

	// Create tblTransactionCash table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "tblTransactionCash" (
			"YQ" TEXT NOT NULL,
			"stockNo"	TEXT NOT NULL,
			"exDividendDate"	TEXT NOT NULL,
			"distributionDate"	TEXT NOT NULL,
			"cashDividend"	REAL NOT NULL,
			"stockDividend"	REAL,
			"quantity"	INTEGER NOT NULL,
			"totalAmount"	INTEGER NOT NULL
		)
	`)
	if err != nil {
		fmt.Println("Error creating tblTransactionCash table:", err)
		return
	}
	fmt.Println("Table tblCashDividend created successfully")

	// Create tblStockMapping table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tblStockMapping (
			stockNo TEXT NOT NULL UNIQUE,
			stockName TEXT NOT NULL,
			PRIMARY KEY(stockNo)
		)
	`)
	if err != nil {
		fmt.Println("Error creating tblStockMapping table:", err)
		return
	}
	fmt.Println("Table tblStockMapping created successfully")

	// Create tblTransaction table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tblTransaction (
			id INTEGER,
			date TEXT NOT NULL,
			time TEXT NOT NULL,
			stockNo TEXT NOT NULL,
			tranType INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			unitPrice REAL NOT NULL,
			totalAmount INTEGER NOT NULL,
			taxes INTEGER NOT NULL,
			PRIMARY KEY("id" AUTOINCREMENT)
		)
	`)
	if err != nil {
		fmt.Println("Error creating tblTransaction table:", err)
		return
	}
	fmt.Println("Table tblTransaction created successfully")

	// Create tblTransactionHistory table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tblTransactionHistory (
			id INTEGER,
			date TEXT NOT NULL,
			time TEXT NOT NULL,
			stockNo TEXT NOT NULL,
			tranType INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			unitPrice REAL NOT NULL,
			totalAmount INTEGER NOT NULL,
			taxes INTEGER NOT NULL,
			PRIMARY KEY("id" AUTOINCREMENT)
		)
	`)
	if err != nil {
		fmt.Println("Error creating tblTransactionHistory table:", err)
		return
	}
	fmt.Println("Table tblTransactionHistory created successfully")

	// Create vv_transactionInventory table
	_, err = db.Exec(`
		CREATE VIEW "vvTransactionInventory" AS
		SELECT 
			stockNo, stockName, tranType, sum(quantity), 
			sum(totalAmount)/sum(quantity) as avgUnitPrice, 
			sum(totalAmount), sum(taxes)
		FROM (
			SELECT a.*, b.stockName
			FROM tblTransaction	a
			JOIN tblStockMapping b on a.stockNo = b.stockNo
		)
		GROUP by stockNo
	`)
	if err != nil {
		fmt.Println("Error creating vv_transactionInventory table:", err)
		return
	}
	fmt.Println("Table vv_transactionInventory created successfully")

	// Create vvTransactionCash table
	_, err = db.Exec(`
		CREATE VIEW "vvTransactionCash" AS 
		SELECT 
			YQ, stockNo, stockName, distributionDate, 
			cashDividend, quantity, totalAmount
		FROM (
			SELECT a.*, b.stockName
			FROM tblTransactionCash	a
			JOIN tblStockMapping b on a.stockNo = b.stockNo
		)
		ORDER BY stockNo, distributionDate
	`)
	if err != nil {
		fmt.Println("Error creating vvTransactionCash table:", err)
		return
	}
	fmt.Println("Table vvTransactionCash created successfully")

	fmt.Println("Database created successfully")
}
